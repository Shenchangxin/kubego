package k8s

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/watch"
	"kubego/global"
	pod_req "kubego/model/pod/request"
	"kubego/response"
	"net/http"
	"strings"
)

type PodApi struct {
}

//因为update的字段属性有限，而我们实际更新过程中会修改定任一字段
func UpdatePod(ctx context.Context, pod *corev1.Pod) error {
	_, err := global.KubeConfigSet.CoreV1().Pods(pod.Namespace).Update(ctx, pod, metav1.UpdateOptions{})
	return err
}

func PatchPod(patchData map[string]interface{}, k8sPod *corev1.Pod, ctx context.Context) error {
	patchDataBytes, _ := json.Marshal(&patchData)

	_, err := global.KubeConfigSet.CoreV1().Pods(k8sPod.Namespace).Patch(
		ctx,
		k8sPod.Name,
		types.StrategicMergePatchType,
		patchDataBytes,
		metav1.PatchOptions{},
	)
	return err
}

func (*PodApi) CreateOrUpdatePod(c *gin.Context) {
	var podReq pod_req.Pod

	if err := c.ShouldBind(&podReq); err != nil {
		response.FailWithMessage(c, "参数解析失败，detail："+err.Error())
		return
	}
	//校验必填项
	if err := podValidate.Validate(&podReq); err != nil {
		response.FailWithMessage(c, "参数验证失败，detail："+err.Error())
		return
	}
	k8sPod := podConvert.PodReq2K8s(podReq)
	ctx := context.TODO()
	podApi := global.KubeConfigSet.CoreV1().Pods(k8sPod.Namespace)
	if k8sGetPod, err := podApi.Get(ctx, k8sPod.Name, metav1.GetOptions{}); err == nil {
		//UpdatePod(ctx, getPod)
		//校验参数是否合理
		k8sPodCopy := *k8sPod
		k8sPodCopy.Name = k8sPod.Name + "-validate"
		_, err := podApi.Create(ctx, &k8sPodCopy, metav1.CreateOptions{
			DryRun: []string{metav1.DryRunAll},
		})
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		//删除
		err = podApi.Delete(ctx, k8sPod.Name, metav1.DeleteOptions{})
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		var labelSelector []string
		for k, v := range k8sGetPod.Labels {
			labelSelector = append(labelSelector, fmt.Sprintf("%s-%s", k, v))
		}
		watcher, err := podApi.Watch(ctx, metav1.ListOptions{
			LabelSelector: strings.Join(labelSelector, ","),
		})
		if err != nil {
			response.FailWithMessage(c, err.Error())
			return
		}
		for event := range watcher.ResultChan() {
			k8sPodChan := event.Object.(*corev1.Pod)

			if _, err := podApi.Get(ctx, k8sPod.Name, metav1.GetOptions{}); k8serror.IsNotFound(err) {
				//重新创建
				if createdPod, err := podApi.Create(ctx, k8sPod, metav1.CreateOptions{}); err != nil {
					errMsg := fmt.Sprintf("Pod[%s-%s]创建失败，detail：%s", k8sPod.Namespace, k8sPod.Name, err.Error())
					response.FailWithMessage(c, errMsg)
					return
				} else {
					successMsg := fmt.Sprintf("Pod[namespace=%s,name=%s]创建成功", createdPod.Namespace, createdPod.Name)
					response.SuccessWithMessage(c, successMsg)
				}
			}
			switch event.Type {
			case watch.Deleted:
				if k8sPodChan.Name != k8sPod.Name {
					continue
				}
				//重新创建
				if createdPod, err := podApi.Create(ctx, k8sPod, metav1.CreateOptions{}); err != nil {
					errMsg := fmt.Sprintf("Pod[%s-%s]创建失败，detail：%s", k8sPod.Namespace, k8sPod.Name, err.Error())
					response.FailWithMessage(c, errMsg)
					return
				} else {
					successMsg := fmt.Sprintf("Pod[namespace=%s,name=%s]创建成功", createdPod.Namespace, createdPod.Name)
					response.SuccessWithMessage(c, successMsg)
				}
			}
		}
	} else {
		createdPod, err := podApi.Create(ctx, k8sPod, metav1.CreateOptions{})
		if err != nil {
			errMsg := fmt.Sprintf("Pod[%s-%s]创建失败，detail：%s", k8sPod.Namespace, k8sPod.Name, err.Error())
			response.FailWithMessage(c, errMsg)
			return
		} else {
			successMsg := fmt.Sprintf("Pod[namespace=%s,name=%s]创建成功", createdPod.Namespace, createdPod.Name)
			response.SuccessWithMessage(c, successMsg)
		}
	}

}

func (*PodApi) GetPodList(c *gin.Context) {

	ctx := context.TODO()
	list, err := global.KubeConfigSet.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Println(err)
	}
	for _, item := range list.Items {
		fmt.Println(item.Namespace, item.Name)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
