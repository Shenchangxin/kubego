package node

import (
	corev1 "k8s.io/api/core/v1"
	"kubego/model/base"
	node_res "kubego/model/node/response"
)

type NodeK8s2Res struct {
}

func getNodeStatus(nodeConditions []corev1.NodeCondition) string {
	nodeStatus := "NotReady"
	for _, condition := range nodeConditions {
		if condition.Type == "Ready" && condition.Status == "True" {
			nodeStatus = "Ready"
			break
		}
	}
	return nodeStatus
}

func getNodeIp(addresses []corev1.NodeAddress, addressType corev1.NodeAddressType) string {
	for _, item := range addresses {
		if item.Type == addressType {
			return item.Address
		}
	}
	return "<none>"
}

func mapToList(m map[string]string) []base.ListMapItem {
	res := make([]base.ListMapItem, 0)
	for k, v := range m {
		res = append(res, base.ListMapItem{
			Key:   k,
			Value: v,
		})
	}
	return res
}

func (this *NodeK8s2Res) GetNodeDetail(nodeK8s corev1.Node) node_res.Node {
	nodeRes := this.GetNodeResItem(nodeK8s)
	//计算label 和 taint
	nodeRes.Taints = nodeK8s.Spec.Taints
	nodeRes.Labels = mapToList(nodeK8s.Labels)
	return nodeRes
}
func (*NodeK8s2Res) GetNodeResItem(nodeK8s corev1.Node) node_res.Node {
	nodeInfo := nodeK8s.Status.NodeInfo
	return node_res.Node{
		Name:             nodeK8s.Name,
		Status:           getNodeStatus(nodeK8s.Status.Conditions),
		Age:              nodeK8s.CreationTimestamp.Unix(),
		InternalIp:       getNodeIp(nodeK8s.Status.Addresses, corev1.NodeInternalIP),
		ExternalIp:       getNodeIp(nodeK8s.Status.Addresses, corev1.NodeExternalIP),
		OsImage:          nodeInfo.OSImage,
		Version:          nodeInfo.KubeletVersion,
		KernelVersion:    nodeInfo.KernelVersion,
		ContainerRuntime: nodeInfo.ContainerRuntimeVersion,
	}
}
