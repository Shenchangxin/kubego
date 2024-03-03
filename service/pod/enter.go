package pod

import "kubego/convert"

type PodServiceGroup struct {
	PodService
}

var podConvert = convert.ConvertGroupApp.PodConvert
