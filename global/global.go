package global

import (
	"k8s.io/client-go/kubernetes"
	"kubego/config"
)

var (
	CONF          config.Server
	KubeConfigSet *kubernetes.Clientset
)
