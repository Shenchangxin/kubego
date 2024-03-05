package k8s

import (
	"kubego/service"
	"kubego/validate"
)

type ApiGroup struct {
	PodApi
	NamespaceApi
	NodeApi
}

var podValidate = validate.ValidateGroupApp.PodValidate
var podService = service.ServiceGroupApp.PodServiceGroup.PodService
var nodeService = service.ServiceGroupApp.NodeServiceGroup
