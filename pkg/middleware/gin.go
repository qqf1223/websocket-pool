package middleware

import (
	"fmt"
	"websocket-pool/global"

	"github.com/gin-gonic/gin"
)

func CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, _ := c.GetQuery(global.CtxKeyXToken)
		fmt.Println(token)
	}
}
