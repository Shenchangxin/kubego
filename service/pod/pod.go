package pod

import (
	"context"
	"errors"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"kubego/global"
	pod_req "kubego/model/pod/request"
	"strings"
)

type PodService struct {
}

func (*PodService) GetPodDetail(namespace string, name string) (podReq pod_req.Pod, err error) {
	ctx := context.TODO()
	podApi := global.KubeConfigSet.CoreV1().Pods(namespace)
	k8sGetPod, err := podApi.Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		errMsg := fmt.Sprintf("Pod[%s-%s]查询失败，detail：%s", namespace, name, err.Error())
		err = errors.New(errMsg)
		return
	}
	podReq = podConvert.K8s2ReqConvert.PodK8s2Req(*k8sGetPod)
	return
}

func (*PodService) CreateOrUpdate(podReq pod_req.Pod) (msg string, err error) {
	k8sPod := podConvert.Req2K8sConvert.PodReq2K8s(podReq)
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
			errMsg := fmt.Sprintf("Pod[%s-%s]创建失败，detail：%s", k8sPod.Namespace, k8sPod.Name, err.Error())

			return errMsg, err
		}
		//删除
		err = podApi.Delete(ctx, k8sPod.Name, metav1.DeleteOptions{})
		if err != nil {
			errMsg := fmt.Sprintf("Pod[%s-%s]创建失败，detail：%s", k8sPod.Namespace, k8sPod.Name, err.Error())
			return errMsg, err
		}
		var labelSelector []string
		for k, v := range k8sGetPod.Labels {
			labelSelector = append(labelSelector, fmt.Sprintf("%s-%s", k, v))
		}
		watcher, err := podApi.Watch(ctx, metav1.ListOptions{
			LabelSelector: strings.Join(labelSelector, ","),
		})
		if err != nil {
			errMsg := fmt.Sprintf("Pod[%s-%s]创建失败，detail：%s", k8sPod.Namespace, k8sPod.Name, err.Error())

			return errMsg, err
		}
		for event := range watcher.ResultChan() {
			k8sPodChan := event.Object.(*corev1.Pod)

			if _, err := podApi.Get(ctx, k8sPod.Name, metav1.GetOptions{}); k8serror.IsNotFound(err) {
				//重新创建
				if createdPod, err := podApi.Create(ctx, k8sPod, metav1.CreateOptions{}); err != nil {
					errMsg := fmt.Sprintf("Pod[%s-%s]创建失败，detail：%s", k8sPod.Namespace, k8sPod.Name, err.Error())
					return errMsg, err
				} else {
					successMsg := fmt.Sprintf("Pod[namespace=%s,name=%s]创建成功", createdPod.Namespace, createdPod.Name)
					return successMsg, err
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
					return errMsg, err
				} else {
					successMsg := fmt.Sprintf("Pod[namespace=%s,name=%s]创建成功", createdPod.Namespace, createdPod.Name)
					return successMsg, err
				}
			}
		}
		return "", nil
	} else {
		if createdPod, err := podApi.Create(ctx, k8sPod, metav1.CreateOptions{}); err != nil {
			errMsg := fmt.Sprintf("Pod[%s-%s]创建失败，detail：%s", k8sPod.Namespace, k8sPod.Name, err.Error())
			return errMsg, err
		} else {
			successMsg := fmt.Sprintf("Pod[namespace=%s,name=%s]创建成功", createdPod.Namespace, createdPod.Name)
			return successMsg, err
		}
	}
}
