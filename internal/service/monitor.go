package service

import (
	"context"
	"websocket-pool/global"
	"websocket-pool/pkg/common/response"
	"websocket-pool/pkg/gredis"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func Monitor(c *gin.Context) {
	// 获取当前客户端连接数

	// 获取请求耗时数量

	// 获取请求失败数量

}

// 获取房间列表
func RoomList(c *gin.Context) {
	res, err := gredis.DoWithContext(context.Background(), "SMEMBERS", global.Redis_Room_List)
	if err != nil {
		return
	}
	roomList := cast.ToStringSlice(res)
	response.OkWithData(roomList, c)
}
