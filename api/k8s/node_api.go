package k8s

import (
	"github.com/gin-gonic/gin"
	node_req "kubego/model/node/request"
	"kubego/response"
)

type NodeApi struct {
}

func (*NodeApi) UpdateNodeTaint(ctx *gin.Context) {
	var updatedTaint node_req.UpdatedTaint
	err := ctx.ShouldBind(&updatedTaint)
	if err != nil {
		response.FailWithMessage(ctx, "参数解析报错！")
		return
	}
	err = nodeService.UpdateNodeTaint(updatedTaint)
	if err != nil {
		response.FailWithMessage(ctx, "更新节点污点(Taint)报错，detail:"+err.Error())
	} else {
		response.Success(ctx)
	}
}

func (*NodeApi) UpdateNodeLabel(ctx *gin.Context) {
	var updatedLabel node_req.UpdatedLabel
	err := ctx.ShouldBind(&updatedLabel)
	if err != nil {
		response.FailWithMessage(ctx, "参数解析报错！")
		return
	}
	err = nodeService.UpdateNodeLabel(updatedLabel)
	if err != nil {
		response.FailWithMessage(ctx, "更新节点标签报错！")
	} else {
		response.Success(ctx)
	}
}

func (*NodeApi) GetNodeDetailOrList(ctx *gin.Context) {
	keyword := ctx.Query("keyword")
	nodeName := ctx.Query("nodeName")
	if nodeName != "" {
		detail, err := nodeService.GetNodeDetail(nodeName)
		if err != nil {
			response.FailWithMessage(ctx, "查询Node详情失败！")
		} else {
			response.SuccessWithDetailed(ctx, "查询Node详情成功！", detail)
		}
	} else {
		list, err := nodeService.GetNodeList(keyword)
		if err != nil {
			response.FailWithMessage(ctx, "查询Node列表失败！")
		} else {
			response.SuccessWithDetailed(ctx, "查询Node列表成功！", list)
		}
	}

}
