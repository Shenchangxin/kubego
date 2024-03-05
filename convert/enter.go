package convert

import (
	"kubego/convert/node"
	"kubego/convert/pod"
)

type ConvertGroup struct {
	PodConvert  pod.PodConvertGroup
	NodeConvert node.Group
}

var ConvertGroupApp = new(ConvertGroup)
