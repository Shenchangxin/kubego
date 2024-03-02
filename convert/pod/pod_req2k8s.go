package pod

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	pod_req "kubego/model/pod/request"
	"strconv"
	"strings"
)

const (
	probe_http = "http"
	probe_tcp  = "tcp"
	probe_exec = "exec"
)
const (
	volume_emptyDir = "emptyDir"
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
			Containers:     pc.getK8sContainers(podReq.Containers),
			InitContainers: pc.getK8sContainers(podReq.InitContainers),
			Volumes:        pc.getK8sVolumes(podReq.Volumes),
			DNSConfig: &corev1.PodDNSConfig{
				Nameservers: podReq.NetWorking.DnsConfig.Nameservers,
			},
			DNSPolicy:     corev1.DNSPolicy(podReq.NetWorking.DnsPolicy),
			HostAliases:   pc.getK8sHostAlias(podReq.NetWorking.HostAliases),
			Hostname:      podReq.NetWorking.HostName,
			RestartPolicy: corev1.RestartPolicy(podReq.Base.RestartPolicy),
		},
	}
}

func (pc *PodConvert) getK8sHostAlias(podReqHostAlias []pod_req.ListMapItem) []corev1.HostAlias {
	podK8sHostAlias := make([]corev1.HostAlias, 0)
	for _, item := range podReqHostAlias {
		podK8sHostAlias = append(podK8sHostAlias, corev1.HostAlias{
			IP:        item.Key,
			Hostnames: strings.Split(item.Value, ","),
		})
	}
	return podK8sHostAlias
}

func (pc *PodConvert) getK8sVolumes(podReqVolumes []pod_req.Volume) []corev1.Volume {
	podK8sVolumes := make([]corev1.Volume, 0)
	for _, volume := range podReqVolumes {
		if volume.Type != volume_emptyDir {
			continue
		}
		source := corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		}
		podK8sVolumes = append(podK8sVolumes, corev1.Volume{
			VolumeSource: source,
			Name:         volume.Name,
		})
	}
	return podK8sVolumes
}

func (pc *PodConvert) getK8sContainers(podReqContainers []pod_req.Container) []corev1.Container {
	podK8sContainers := make([]corev1.Container, 0)
	for _, item := range podReqContainers {

		podK8sContainers = append(podK8sContainers, pc.getK8sContainer(item))

	}
	return podK8sContainers
}

func (pc *PodConvert) getK8sContainer(podReqContainer pod_req.Container) corev1.Container {
	return corev1.Container{
		Name:            podReqContainer.Name,
		Image:           podReqContainer.Image,
		ImagePullPolicy: corev1.PullPolicy(podReqContainer.ImagePullPolicy),
		TTY:             podReqContainer.Tty,
		Command:         podReqContainer.Command,
		Args:            podReqContainer.Args,
		WorkingDir:      podReqContainer.WorkingDir,
		SecurityContext: &corev1.SecurityContext{
			Privileged: &podReqContainer.Privileged,
		},
		Ports:          pc.getK8sPorts(podReqContainer.Ports),
		Env:            pc.getK8sEnv(podReqContainer.Envs),
		VolumeMounts:   pc.getK8sVolumeMount(podReqContainer.VolumeMounts),
		StartupProbe:   pc.getK8sContainerProbe(podReqContainer.StartupProbe),
		LivenessProbe:  pc.getK8sContainerProbe(podReqContainer.LivenessProbe),
		ReadinessProbe: pc.getK8sContainerProbe(podReqContainer.ReadinessProbe),
		Resources:      pc.getK8sResources(podReqContainer.Resources),
	}
}

func (pc *PodConvert) getK8sPorts(podReqPorts []pod_req.ContainerPort) []corev1.ContainerPort {
	podK8sContainerPorts := make([]corev1.ContainerPort, 0)
	for _, item := range podReqPorts {
		podK8sContainerPorts = append(podK8sContainerPorts, corev1.ContainerPort{
			Name:          item.Name,
			HostPort:      item.HostPort,
			ContainerPort: item.ContainerPort,
		})
	}
	return podK8sContainerPorts
}

func (pc *PodConvert) getK8sResources(podReqResources pod_req.Resources) corev1.ResourceRequirements {
	var k8sPodResources corev1.ResourceRequirements
	if !podReqResources.Enable {
		return k8sPodResources
	}
	k8sPodResources.Requests = corev1.ResourceList{
		corev1.ResourceCPU:    resource.MustParse(strconv.Itoa(int(podReqResources.CpuRequest)) + "m"),
		corev1.ResourceMemory: resource.MustParse(strconv.Itoa(int(podReqResources.MemRequest)) + "Mi"),
	}
	k8sPodResources.Limits = corev1.ResourceList{
		corev1.ResourceCPU:    resource.MustParse(strconv.Itoa(int(podReqResources.CpuLimit)) + "m"),
		corev1.ResourceMemory: resource.MustParse(strconv.Itoa(int(podReqResources.MemLimit)) + "Mi"),
	}

	return k8sPodResources
}

func (pc *PodConvert) getK8sContainerProbe(podReqProbe pod_req.ContainerProbe) *corev1.Probe {
	if !podReqProbe.Enable {
		return nil
	}
	var k8sProbe corev1.Probe
	switch podReqProbe.Type {
	case probe_http:
		httpGet := podReqProbe.HttpGet
		k8sHttpHeaders := make([]corev1.HTTPHeader, 0)
		for _, header := range httpGet.HttpHeaders {
			k8sHttpHeaders = append(k8sHttpHeaders, corev1.HTTPHeader{
				Name:  header.Key,
				Value: header.Value,
			})
		}
		k8sProbe.HTTPGet = &corev1.HTTPGetAction{
			Scheme:      corev1.URIScheme(httpGet.Scheme),
			Host:        httpGet.Host,
			Port:        intstr.FromInt(int(httpGet.Port)),
			Path:        httpGet.Path,
			HTTPHeaders: k8sHttpHeaders,
		}
	case probe_exec:
		exec := podReqProbe.Exec
		k8sProbe.Exec = &corev1.ExecAction{
			Command: exec.Command,
		}
	case probe_tcp:
		tcpSocket := podReqProbe.TcpSocket

		k8sProbe.TCPSocket = &corev1.TCPSocketAction{
			Host: tcpSocket.Host,
			Port: intstr.FromInt(int(tcpSocket.Port)),
		}

	}
	return &k8sProbe
}

func (pc *PodConvert) getK8sVolumeMount(podReqMounts []pod_req.VolumeMount) []corev1.VolumeMount {
	podK8sVolumeMounts := make([]corev1.VolumeMount, 0)
	for _, mount := range podReqMounts {
		podK8sVolumeMounts = append(podK8sVolumeMounts, corev1.VolumeMount{
			Name:      mount.MountName,
			MountPath: mount.MountPath,
			ReadOnly:  mount.ReadOnly,
		})
	}
	return podK8sVolumeMounts
}

func (pc *PodConvert) getK8sEnv(podReqEnv []pod_req.ListMapItem) []corev1.EnvVar {
	podK8sEnvs := make([]corev1.EnvVar, 0)
	for _, item := range podReqEnv {
		podK8sEnvs = append(podK8sEnvs, corev1.EnvVar{
			Name:  item.Key,
			Value: item.Value,
		})
	}
	return podK8sEnvs
}

func (pc *PodConvert) getK8sLabels(podReqLabels []pod_req.ListMapItem) map[string]string {
	podK8sLabels := make(map[string]string)
	for _, label := range podReqLabels {
		podK8sLabels[label.Key] = label.Value
	}
	return podK8sLabels
}
