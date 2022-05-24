package middlewares

import (
	"gitee.com/tzxhy/web/constants"
	"gitee.com/tzxhy/web/models"
	"gitee.com/tzxhy/web/utils"
	"github.com/gin-gonic/gin"
)

func MayAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		token, err := c.Cookie(constants.TOKEN_COOKIE_NAME)
		if err != nil {
			c.Next()
			return
		}

		mc, err := utils.ParseToken(token)

		if err != nil {
			c.Next()
			return
		}

		has := models.GetKey(token)
		if has == nil {
			c.Next()
			return
		}

		// 将当前请求的username信息保存到请求的上下文c上
		c.Set("username", mc.Username)
		c.Set("uid", mc.UserId)
		c.Set("token", token)
		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
}
