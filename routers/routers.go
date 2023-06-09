package routers

import (
	"websocket-pool/internal/service"
	"websocket-pool/pkg/server"

	"github.com/gin-gonic/gin"
)

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
	r.GET("/wspool", server.Ws.WebsocketEntry)
	r.POST("/wspool/client/message", service.ServerEntry)
	return r
}
