package main

import (
	"websocket-pool/protobuf"
	"websocket-pool/routers"

	"websocket-pool/pkg/config"
	"websocket-pool/pkg/gredis"
	"websocket-pool/pkg/log"
	"websocket-pool/pkg/server"

	"websocket-pool/global"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
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
	// websocket server init
	srvs.BindServer(server.Ws.Init(engine))
	// rpc server init
	srvs.BindServer(&server.GRPCServer{
		Name:        "Tool",
		Addr:        global.GVA_CONFIG.Rpc.Host,
		ServiceOpts: []grpc.ServerOption{},
		RegisterFn: func(srv *grpc.Server) {
			protobuf.RegisterWsServerServer(srv, nil)
		},
	})
	// srvs.BindServer(server.Rpc.Init(engine))
	srvs.Run()
}

func initConfig() {
	// 初始化viper
	global.GVA_VP = config.Viper()
	config.Viper("./conf/config.yaml")
}

func initRedis() error {
	gredis.Init(gredis.NewConfig())
	return nil
}

func initLog() {
	// 初始化zap
	global.GVA_LOG = log.InitLogger()
	zap.ReplaceGlobals(global.GVA_LOG)
}

// go generate
