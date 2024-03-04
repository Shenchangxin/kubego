package pod

import (
	corev1 "k8s.io/api/core/v1"
	pod_req "kubego/model/pod/request"
	"strings"
)

const VOLUME_TYPE_EMPTYDIR = "emptyDir"

type K8s2ReqConvert struct {
}

func (*K8s2ReqConvert) PodK8s2Req(podK8s corev1.Pod) pod_req.Pod {
	return pod_req.Pod{
		Base:           getReqBase(podK8s),
		NetWorking:     getReqNetWorking(podK8s),
		Volumes:        getReqVolumes(podK8s.Spec.Volumes),
		Containers:     nil,
		InitContainers: nil,
	}
}

func getReqContainers(containersK8s []corev1.Container) []pod_req.Container {
	podReqContainers := make([]pod_req.Container, 0)
	for _, item := range containersK8s {
		reqContainer := getReqContainer(item)
		podReqContainers = append(podReqContainers, reqContainer)
	}
	return podReqContainers
}
func getReqContainer(container corev1.Container) pod_req.Container {
	return pod_req.Container{
		Name:            container.Name,
		Image:           container.Image,
		ImagePullPolicy: string(container.ImagePullPolicy),
		Tty:             container.TTY,
		WorkingDir:      container.WorkingDir,
		Command:         container.Command,
		Args:            container.Args,
		Ports:           getReqContainerPorts(container.Ports),

		Envs:           getReqContainerEnvs(container.Env),
		Privileged:     "",
		Resources:      nil,
		VolumeMounts:   nil,
		StartupProbe:   nil,
		LivenessProbe:  nil,
		ReadinessProbe: nil,
	}
}

func getReqContainerEnvs(envsK8s []corev1.EnvVar) []pod_req.ListMapItem {
	envsReq := make([]pod_req.ListMapItem, 0)
	for _, item := range envsK8s {
		envsReq = append(envsReq, pod_req.ListMapItem{
			Key:   item.Name,
			Value: item.Value,
		})
	}
	return envsReq
}

func getReqContainerPorts(portsK8s []corev1.ContainerPort) []pod_req.ContainerPort {
	portsReq := make([]pod_req.ContainerPort, 0)
	for _, item := range portsK8s {
		portsReq = append(portsReq, pod_req.ContainerPort{
			Name:          item.Name,
			HostPort:      item.HostPort,
			ContainerPort: item.ContainerPort,
		})
	}
	return portsReq
}

func getReqVolumes(volumes []corev1.Volume) []pod_req.Volume {
	volumesReq := make([]pod_req.Volume, 0)
	for _, volume := range volumes {
		if volume.EmptyDir == nil {
			continue
		}
		volumesReq = append(volumesReq, pod_req.Volume{
			Type: VOLUME_TYPE_EMPTYDIR,
			Name: volume.Name,
		})
	}
	return volumesReq
}

func getReqHostAliases(hostAlias []corev1.HostAlias) []pod_req.ListMapItem {

	hostAliasReq := make([]pod_req.ListMapItem, 0)
	for _, alias := range hostAlias {
		hostAliasReq = append(hostAliasReq, pod_req.ListMapItem{
			Key:   alias.IP,
			Value: strings.Join(alias.Hostnames, ","),
		})
	}
	return hostAliasReq
}

func getReqDnsConfig(dnsConfigK8s *corev1.PodDNSConfig) pod_req.DnsConfig {
	var dnsConfigReq pod_req.DnsConfig
	if dnsConfigK8s != nil {
		dnsConfigReq.Nameservers = dnsConfigK8s.Nameservers
	}

	return dnsConfigReq

}

func getReqNetWorking(pod corev1.Pod) pod_req.NetWorking {
	return pod_req.NetWorking{
		HostNetwork: pod.Spec.HostNetwork,
		HostName:    pod.Spec.Hostname,
		DnsPolicy:   string(pod.Spec.DNSPolicy),
		DnsConfig:   getReqDnsConfig(pod.Spec.DNSConfig),
		HostAliases: getReqHostAliases(pod.Spec.HostAliases),
	}
}
func getReqLabels(data map[string]string) []pod_req.ListMapItem {
	labels := make([]pod_req.ListMapItem, 0)
	for k, v := range data {
		labels = append(labels, pod_req.ListMapItem{
			Key:   k,
			Value: v,
		})
	}
	return labels
}

func getReqBase(pod corev1.Pod) pod_req.Base {
	return pod_req.Base{
		Name:          pod.Name,
		Namespace:     pod.Namespace,
		Labels:        getReqLabels(pod.Labels),
		RestartPolicy: string(pod.Spec.RestartPolicy),
	}
}
