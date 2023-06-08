package main

import (
	"fmt"
	"net"
	"strings"
	"websocket-pool/entity"
	"websocket-pool/internal/rpc"
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
	r := routers.Init()
	// websocket server init
	srvs.BindServer(server.Ws.Init(r))
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
	initNodeInfo()
	srvs.Run()
}
func initNodeInfo() (entity.MessageBody, error) {
	ret := entity.MessageBody{}
	conn, err := net.Dial("udp", "192.168.2.3:80")
	if err != nil {
		return ret, err
	}
	defer conn.Close()
	fmt.Println("当前ip:", strings.Split(conn.LocalAddr().String(), ":")[0])
	return ret, nil
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
