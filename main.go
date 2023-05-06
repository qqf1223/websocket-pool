package main

import (
	"fmt"
	"websocket-pool/routers"

	"websocket-pool/pkg/config"
	"websocket-pool/pkg/gredis"
	"websocket-pool/pkg/log"
	"websocket-pool/pkg/server"

	"websocket-pool/global"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 服务初始化
	srvs := server.NewServers()
	// 初始化配置文件
	initConfig()

	// 初始化日志组件
	initLog()

	srvs.BindBeforeHandler(func() error {

		// 初始化redis
		initRedis()
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
	// 初始化viper
	global.GVA_VP = config.Viper()
	config.Viper("./conf/config.yaml")
}

func initRedis() error {
	gredis.Init(gredis.NewConfig())
	fmt.Printf("redis pool init success\n")
	return nil
}

func initLog() {
	// 初始化zap
	global.GVA_LOG = log.InitLogger()
	zap.ReplaceGlobals(global.GVA_LOG)
}
