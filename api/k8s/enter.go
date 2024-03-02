package k8s

import (
	"kubego/convert"
	"kubego/validate"
)

type ApiGroup struct {
	PodApi
	NamespaceApi
}

var podValidate = validate.ValidateGroupApp.PodValidate
var podConvert = convert.ConvertGroupApp.PodConvert
