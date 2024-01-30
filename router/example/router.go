package example

import (
	"github.com/gin-gonic/gin"
	"kubego/api"
)

type ExampleRouter struct {
}

func (*ExampleRouter) InitExample(r *gin.Engine) {
	group := r.Group("/example")
	exampleGroup := api.ApiGroupApp.ExampleApiGroup
	group.GET("/ping", exampleGroup.ExampleTest)
}
