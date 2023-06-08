package service

import (
	"websocket-pool/entity"
	"websocket-pool/pkg/client"

	"github.com/gin-gonic/gin"
)

type ServerService struct {
}

func ServerEntry(c *gin.Context) {
	req := &entity.BizMessageEntity{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return
	}

	client.NewClientManager().SendTo(c, req)
}
