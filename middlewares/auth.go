package middlewares

import (
	"net/http"

	"gitee.com/tzxhy/web/constants"
	"gitee.com/tzxhy/web/models"
	"gitee.com/tzxhy/web/utils"
	"github.com/gin-gonic/gin"
)

func NeedAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		token, err := c.Cookie(constants.TOKEN_COOKIE_NAME)
		if err != nil {
			c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_NOT_LOGIN_TIPS.Code, constants.CODE_NOT_LOGIN_TIPS.Tip, nil))
			c.Abort()
			return
		}

		mc, err := utils.ParseToken(token)

		if err != nil {
			c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_TOKEN_NOT_VALID_TIPS.Code, constants.CODE_TOKEN_NOT_VALID_TIPS.Tip, nil))
			c.Abort()
			return
		}

		has := models.GetKey(token)
		if has == nil {
			c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_TOKEN_NOT_VALID_TIPS.Code, constants.TIPS_TOKEN_VALID_WITH_ERROR, nil))
			c.Abort()
			return
		}

		// 将当前请求的username信息保存到请求的上下文c上
		c.Set("username", mc.Username)
		c.Set("uid", mc.UserId)
		c.Set("token", token)
		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
}
