package k8s

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubego/global"
	"kubego/response"
	"net/http"
)

type PodApi struct {
}

func (*PodApi) CreateOrUpdatePod(c *gin.Context) {
	response.Success(c)
}

func (*PodApi) GetPodList(c *gin.Context) {

	ctx := context.TODO()
	list, err := global.KubeConfigSet.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Println(err)
	}
	for _, item := range list.Items {
		fmt.Println(item.Namespace, item.Name)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
