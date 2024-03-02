package convert

import "kubego/convert/pod"

type ConvertGroup struct {
	PodConvert pod.PodConvert
}

var ConvertGroupApp = new(ConvertGroup)
