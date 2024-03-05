package node

import "kubego/convert"

type Group struct {
	NodeService
}

var nodeConvert = convert.ConvertGroupApp.NodeConvert
