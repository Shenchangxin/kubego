package initialize

import (
	"github.com/gin-gonic/gin"
	"kubego/router"
)

func Routers() *gin.Engine {
	r := gin.Default()
	exampleGroup := router.RouterGroupApp.ExampleRouterGroup
	k8SGroup := router.RouterGroupApp.K8SRouterGroup
	exampleGroup.InitExample(r)
	k8SGroup.InitK8SRouter(r)
	return r
}
