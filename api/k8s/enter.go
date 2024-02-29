package k8s

import "kubego/validate"

type ApiGroup struct {
	PodApi
	NamespaceApi
}

var podValidate = validate.ValidateGroupApp.PodValidate
