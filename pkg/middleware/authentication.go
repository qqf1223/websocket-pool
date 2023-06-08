package middleware

import (
	"fmt"
	"websocket-pool/global"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		var success bool
		token, _ := c.GetQuery(global.CtxKeyXToken)
		fmt.Println(token)
		// grpcclient.Authentication(token)
		if success {

			c.Next()
		} else {
			c.Abort()
			return
		}
	}
}
