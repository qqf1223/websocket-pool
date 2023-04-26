package routers

import (
	"websocket-pool/pkg/server"

	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	r.GET("/ws", server.Ws.WebsocketEntry)
}
