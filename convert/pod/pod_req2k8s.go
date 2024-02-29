package pod

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	pod_req "kubego/model/pod/request"
)

type PodConvert struct {
}

//将Pod的请求数据转换为K8s结构的数据
func (pc *PodConvert) PodReq2K8s(podReq pod_req.Pod) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podReq.Base.Name,
			Namespace: podReq.Base.Namespace,
			Labels:    pc.getK8sLabels(podReq.Base.Labels),
		},
		Spec: corev1.PodSpec{
			Containers:     nil,
			InitContainers: nil,
			Volumes:        nil,
			DNSConfig:      &corev1.PodDNSConfig{},
			DNSPolicy:      "",
			HostAliases:    nil,
			Hostname:       "",
			RestartPolicy:  "",
		},
	}
}
func (*PodConvert) getK8sLabels(podReqLabels []pod_req.ListMapItem) map[string]string {
	podK8sLabels := make(map[string]string)
	for _, label := range podReqLabels {
		podK8sLabels[label.Key] = label.Value
	}
	return podK8sLabels
}
