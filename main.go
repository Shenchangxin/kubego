package main

import (
	"kubego/global"
	"kubego/initialize"
)

func main() {
	r := initialize.Routers()
	initialize.Viper()
	initialize.K8S()
	panic(r.Run(global.CONF.System.Addr))
}
