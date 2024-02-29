package validate

import (
	"errors"
	pod_req "kubego/model/pod/request"
)

const (
	IMAGE_PULL_POLICY_IFNOTPRESENT = "IfNotPresent"
)

const (
	RESTART_POLICY_ALWAYS = "Always"
)

type PodValidate struct {
}

func (*PodValidate) Validate(podReq *pod_req.Pod) error {

	//校验必填项
	if podReq.Base.Name == "" {
		return errors.New("请定义Pod的名字！")
	}
	if len(podReq.Containers) == 0 {
		return errors.New("请定义Pod的容器信息！")
	}
	if len(podReq.InitContainers) > 0 {
		for index, container := range podReq.InitContainers {
			if container.Name == "" {
				return errors.New("InitContainer中发现没有定义名称的容器！")
			}
			if container.Image == "" {
				return errors.New("InitContainer中发现没有定义镜像的容器！")
			}
			//对非必填项赋值默认值
			if container.ImagePullPolicy == "" {
				podReq.InitContainers[index].ImagePullPolicy = IMAGE_PULL_POLICY_IFNOTPRESENT
			}
		}
	}

	if len(podReq.Containers) > 0 {
		for index, container := range podReq.InitContainers {
			if container.Name == "" {
				return errors.New("Containers中发现没有定义名称的容器！")
			}
			if container.Image == "" {
				return errors.New("Containers中发现没有定义镜像的容器！")
			}
			//对非必填项赋值默认值
			if container.ImagePullPolicy == "" {
				podReq.Containers[index].ImagePullPolicy = IMAGE_PULL_POLICY_IFNOTPRESENT
			}
		}
	}
	if podReq.Base.RestartPolicy == "" {
		podReq.Base.RestartPolicy = RESTART_POLICY_ALWAYS
	}
	return nil
}
