package main

import (
	"strings"
	"websocket-pool/internal/rpc"
	"websocket-pool/pkg/config"
	"websocket-pool/pkg/gredis"
	"websocket-pool/pkg/log"
	"websocket-pool/pkg/server"
	"websocket-pool/pkg/utils"
	"websocket-pool/protobuf"
	"websocket-pool/routers"

	"websocket-pool/global"
	"websocket-pool/pkg/discover/getcdv3"

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

	// websocket server init
	srvs.BindServer(server.Ws.Init(routers.Init()))

	// rpc server init
	srvs.BindServer(&server.GRPCServer{
		Name:        "RPC",
		Addr:        global.GVA_CONFIG.Rpc.Addr,
		ServiceOpts: []grpc.ServerOption{},
		RegisterFn: func(srv *grpc.Server) {
			protobuf.RegisterWsServerServer(srv, new(rpc.WsService))
		},
	})
	// udp server init
	srvs.BindServer(&server.UDPServer{
		Recv: make(chan []byte, 1000),
	})
	// http server init
	srvs.BindServer(&server.WebServer{
		Name:    "Web",
		Addr:    global.GVA_CONFIG.Http.Addr,
		Timeout: global.GVA_CONFIG.Http.Timeout,
		Handler: routers.WebInit(),
	})
	go registerNodeInfo()
	srvs.Run()
}
func registerNodeInfo() (err error) {
	registerIP, _ := utils.GetOutboundIP()
	err = getcdv3.RegisterEtcd(global.GVA_CONFIG.Etcd.Schema, strings.Join(global.GVA_CONFIG.Etcd.Addr, ","), registerIP, 0, global.GVA_CONFIG.System.Name, 10)
	if err != nil {
		global.GVA_LOG.Error("RegisterEtcd failed ", zap.Error(err))
		return
	}
	global.GVA_LOG.Info("RegisterEtcd ", zap.String("schema", global.GVA_CONFIG.Etcd.Schema), zap.String("addr", strings.Join(global.GVA_CONFIG.Etcd.Addr, ",")), zap.String("registerIP", registerIP), zap.String("srvName", global.GVA_CONFIG.System.Name))
	return
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
