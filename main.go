package main

import (
	"websocket-pool/routers"

	"websocket-pool/pkg/config"
	"websocket-pool/pkg/gredis"
	"websocket-pool/pkg/server"

	"github.com/gin-gonic/gin"
)

func main() {
	// 服务初始化
	srvs := server.NewServers()
	// 初始化配置文件
	initConfig()

	// 初始化日志组件
	initLog()

	// 初始化redis
	srvs.BindBeforeHandler(func() error {
		gredis.Init(gredis.NewConfig())
		return nil
	})

	gin.SetMode(gin.ReleaseMode)
	// 路由初始化
	engine := gin.Default()

	routers.Init(engine)

	srvs.BindServer(server.Ws.Init(engine))
	srvs.Run()
}

func initConfig() {
	config.Viper("./conf/config.yaml")
}

func initRedis() {

}

func initLog() {

}
