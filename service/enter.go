package service

import (
	"kubego/service/node"
	"kubego/service/pod"
)

type ServiceGroup struct {
	PodServiceGroup  pod.PodServiceGroup
	NodeServiceGroup node.Group
}

var ServiceGroupApp = new(ServiceGroup)
