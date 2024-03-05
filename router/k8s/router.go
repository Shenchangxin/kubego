package k8s

import (
	"github.com/gin-gonic/gin"
	"kubego/api"
)

type K8SRouter struct {
}

func (*K8SRouter) InitK8SRouter(r *gin.Engine) {
	group := r.Group("/k8s")
	apiGroup := api.ApiGroupApp.K8SApiGroup
	group.POST("/createOrUpdatePod", apiGroup.CreateOrUpdatePod)
	group.GET("/getPodListOrDetail", apiGroup.GetPodListOrDetail)
	group.GET("/getNameSpace", apiGroup.GetNamespaceList)
	group.DELETE("/deletePod/:namespace/:name", apiGroup.DeletePod)

	group.GET("/node", apiGroup.GetNodeDetailOrList)
	group.PUT("/node/label", apiGroup.UpdateNodeLabel)
	group.PUT("/node/taint", apiGroup.UpdateNodeTaint)
}
