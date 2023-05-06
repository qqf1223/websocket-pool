package routers

import (
	"websocket-pool/pkg/server"

	"websocket-pool/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {
	r.Use(gin.Recovery())
	r.Use(middleware.CheckLogin())
	r.GET("/ws", server.Ws.WebsocketEntry)
}
