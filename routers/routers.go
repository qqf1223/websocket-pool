package routers

import (
	"websocket-pool/internal/service"
	"websocket-pool/pkg/server"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}
func Init() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	// r.Use(middleware.Authentication())
	r.GET("/wspool", server.Ws.WebsocketEntry)
	r.POST("/wspool/client/monitor", service.ServerEntry)
	return r
}

func WebInit() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	// r.Use(middleware.Authentication())
	r.GET("wspool/roomList", service.RoomList) // 获取房间列表
	r.POST("/wspool/monitor", service.Monitor)
	return r
}
