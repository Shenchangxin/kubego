package example

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ExampleApi struct {
}

func (*ExampleApi) ExampleTest(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
