package middlewares

import (
	"net/http"

	"gitee.com/tzxhy/web/utils"
	"github.com/gin-gonic/gin"
)

var ForbiddenIps = [...]string{
	"10.8.1.125",
}

func IpForbidden() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIp := c.ClientIP()
		ForbiddenIpSlice := ForbiddenIps[:]
		if utils.Has(&ForbiddenIpSlice, clientIp) {
			c.String(http.StatusForbidden, "Forbidden")
			c.Abort()
			return
		}
		c.Next()
	}
}
