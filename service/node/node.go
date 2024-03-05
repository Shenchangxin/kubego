package node

import (
	"context"
	"encoding/json"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"kubego/global"
	node_req "kubego/model/node/request"
	node_res "kubego/model/node/response"
	"strings"
)

type NodeService struct {
}

func (*NodeService) UpdateNodeTaint(updatedTaint node_req.UpdatedTaint) error {
	patchData := map[string]any{
		"spec": map[string]any{
			"taints": updatedTaint.Taints,
		},
	}
	patchDataBytes, _ := json.Marshal(&patchData)
	_, err := global.KubeConfigSet.CoreV1().Nodes().Patch(
		context.TODO(),
		updatedTaint.Name,
		types.StrategicMergePatchType,
		patchDataBytes,
		metav1.PatchOptions{},
	)
	return err
}

func (*NodeService) UpdateNodeLabel(updatedLabel node_req.UpdatedLabel) error {
	labelsMap := make(map[string]string, 0)
	for _, label := range updatedLabel.Labels {
		labelsMap[label.Key] = label.Value
	}
	labelsMap["$patch"] = "replace"
	patchData := map[string]any{
		"metadata": map[string]any{
			"labels": labelsMap,
		},
	}
	patchDataBytes, _ := json.Marshal(&patchData)
	_, err := global.KubeConfigSet.CoreV1().Nodes().Patch(
		context.TODO(),
		updatedLabel.Name,
		types.StrategicMergePatchType,
		patchDataBytes,
		metav1.PatchOptions{},
	)
	return err
}

func (*NodeService) GetNodeDetail(nodeName string) (*node_res.Node, error) {
	nodeK8s, err := global.KubeConfigSet.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	detail := nodeConvert.GetNodeDetail(*nodeK8s)
	return &detail, err

}
func (*NodeService) GetNodeList(keyword string) ([]node_res.Node, error) {
	ctx := context.TODO()
	list, err := global.KubeConfigSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	nodeResList := make([]node_res.Node, 0)
	for _, item := range list.Items {
		if strings.Contains(item.Name, keyword) {
			nodeRes := nodeConvert.GetNodeResItem(item)
			nodeResList = append(nodeResList, nodeRes)
		}
	}
	return nodeResList, err
}
