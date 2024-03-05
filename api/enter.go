package api

import (
	"kubego/api/example"
	"kubego/api/k8s"
)

type ApiGroup struct {
	ExampleApiGroup example.ApiGroup
	K8SApiGroup     k8s.ApiGroup
	NodeApiGroup    k8s.NodeApi
}

var ApiGroupApp = new(ApiGroup)
