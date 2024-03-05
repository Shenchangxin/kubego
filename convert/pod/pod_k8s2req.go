package pod

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	pod_req "kubego/model/pod/request"
	pod_res "kubego/model/pod/response"
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
		Containers:     getReqContainers(podK8s.Spec.Containers),
		InitContainers: getReqContainers(podK8s.Spec.InitContainers),
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
		Privileged:     getReqContainerPrivileged(container.SecurityContext),
		Resources:      getReqContainerResources(container.Resources),
		VolumeMounts:   getReqContainerVolumeMounts(container.VolumeMounts),
		StartupProbe:   getReqContainerProbe(container.StartupProbe),
		LivenessProbe:  getReqContainerProbe(container.LivenessProbe),
		ReadinessProbe: getReqContainerProbe(container.ReadinessProbe),
	}
}
func getReqContainerProbe(probeK8s *corev1.Probe) pod_req.ContainerProbe {
	containerProbe := pod_req.ContainerProbe{
		Enable: false,
	}
	//先判断是否探针为空
	if probeK8s != nil {
		containerProbe.Enable = true
		//再判断 探针具体是什么类型
		if probeK8s.Exec != nil {
			containerProbe.Type = probe_exec
			containerProbe.Exec.Command = probeK8s.Exec.Command
		} else if probeK8s.HTTPGet != nil {
			containerProbe.Type = probe_http
			httpGet := probeK8s.HTTPGet
			headersReq := make([]pod_req.ListMapItem, 0)
			for _, headerK8s := range httpGet.HTTPHeaders {
				headersReq = append(headersReq, pod_req.ListMapItem{
					Key:   headerK8s.Name,
					Value: headerK8s.Value,
				})
			}
			containerProbe.HttpGet = pod_req.ProbeHttpGet{
				Host:        httpGet.Host,
				Port:        httpGet.Port.IntVal,
				Scheme:      string(httpGet.Scheme),
				Path:        httpGet.Path,
				HttpHeaders: headersReq,
			}
		} else if probeK8s.TCPSocket != nil {
			containerProbe.Type = probe_tcp
			containerProbe.TcpSocket = pod_req.ProbeTcpSocket{
				Host: probeK8s.TCPSocket.Host,
				Port: probeK8s.TCPSocket.Port.IntVal,
			}
		} else {
			containerProbe.Type = probe_http
			return containerProbe
		}
		containerProbe.InitialDelaySeconds = probeK8s.InitialDelaySeconds
		containerProbe.PeriodSeconds = probeK8s.PeriodSeconds
		containerProbe.TimeoutSeconds = probeK8s.TimeoutSeconds
		containerProbe.SuccessThreshold = probeK8s.SuccessThreshold
		containerProbe.FailureThreshold = probeK8s.FailureThreshold
	}
	return containerProbe
}
func (*K8s2ReqConvert) PodK8s2ItemRes(pod corev1.Pod) pod_res.PodListItem {

	var totalC, readyC, restartC int32
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if containerStatus.Ready {
			readyC++
		}
		restartC += containerStatus.RestartCount
		totalC++
	}
	var podStatus string
	if pod.Status.Phase != "Running" {
		podStatus = "Error"
	} else {
		podStatus = "Running"
	}
	return pod_res.PodListItem{
		Name:     pod.Name,
		Ready:    fmt.Sprintf("%d/%d", readyC, totalC),
		Status:   podStatus,
		Restarts: restartC,
		Age:      pod.CreationTimestamp.Unix(),
		IP:       pod.Status.PodIP,
		Node:     pod.Spec.NodeName,
	}
}

func getReqContainerVolumeMounts(volumeMountsK8s []corev1.VolumeMount) []pod_req.VolumeMount {
	volumesReq := make([]pod_req.VolumeMount, 0)
	for _, item := range volumeMountsK8s {

		volumesReq = append(volumesReq, pod_req.VolumeMount{
			MountName: item.Name,
			MountPath: item.MountPath,
			ReadOnly:  item.ReadOnly,
		})
	}
	return volumesReq
}
func getReqContainerResources(requirements corev1.ResourceRequirements) pod_req.Resources {
	reqResources := pod_req.Resources{
		Enable: false,
	}
	requests := requirements.Requests
	limits := requirements.Limits
	if requests != nil {
		reqResources.Enable = true
		reqResources.CpuRequest = int32(requests.Cpu().MilliValue()) // m
		//MiB
		reqResources.MemRequest = int32(requests.Memory().Value() / (1024 * 1024)) //Bytes
	}
	if limits != nil {
		reqResources.Enable = true
		reqResources.CpuLimit = int32(limits.Cpu().MilliValue())
		reqResources.MemLimit = int32(limits.Memory().Value() / (1024 * 1024))
	}
	return reqResources
}

func getReqContainerPrivileged(ctx *corev1.SecurityContext) (privileged bool) {
	if ctx != nil {
		privileged = *ctx.Privileged
	}
	return
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
