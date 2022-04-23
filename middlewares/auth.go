package middlewares

import (
	"log"
	"net/http"

	"gitee.com/tzxhy/web/constants"
	"gitee.com/tzxhy/web/utils"
	"github.com/gin-gonic/gin"
)

func NeedAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 这里的具体实现方式要依据你的实际业务情况决定
		token, err := c.Cookie(constants.TOKEN_COOKIE_NAME)
		if err != nil {
			c.JSON(http.StatusForbidden, utils.ReturnJSON(constants.CODE_NOT_LOGIN, constants.TIPS_NOT_LOGIN, nil))
			c.Abort()
			return
		}

		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := utils.ParseToken(token)
		log.Println(mc)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 2005,
				"msg":  "无效的Token",
			})
			c.Abort()
			return
		}
		// 将当前请求的username信息保存到请求的上下文c上
		c.Set("username", mc.Username)
		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
}
