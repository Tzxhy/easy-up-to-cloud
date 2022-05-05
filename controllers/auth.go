package controllers

import (
	"log"
	"net/http"
	"time"

	"gitee.com/tzxhy/web/constants"
	"gitee.com/tzxhy/web/models"
	"gitee.com/tzxhy/web/utils"
	"github.com/gin-gonic/gin"
)

type LoginInfo struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
	Remember bool   `json:"remember" form:"remember"`
}

// 注册
func Register(c *gin.Context) {
	var loginInfo LoginInfo
	if c.ShouldBind(&loginInfo) == nil {
		log.Print(loginInfo)
		if loginInfo.Username != "" && loginInfo.Password != "" {

			alreadyHasUserName := models.HasUsername(loginInfo.Username)
			if alreadyHasUserName { // 用户名已注册
				c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_USERNAME_IS_REGISTERED, constants.TIPS_USERNAME_IS_REGISTERED, nil))
				return
			}
			_, err := models.AddUser(loginInfo.Username, loginInfo.Password)
			if err == nil {
				c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_OK, "", nil))
				return
			}
		}
	}
	c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_REGISTER_PARAM_NOT_VALID, constants.TIPS_REGISTER_PARAM_NOT_VALID, nil))
}

// 登录
func Login(c *gin.Context) {
	var loginInfo LoginInfo
	if c.ShouldBind(&loginInfo) == nil {
		if loginInfo.Username != "" && loginInfo.Password != "" {
			userInfo := models.GetUserByNameAndPassword(loginInfo.Username, loginInfo.Password)
			if userInfo != nil {
				tokenString, _ := utils.GenToken(userInfo.Username, userInfo.Uid)
				timeSecond := 0
				if loginInfo.Remember {
					timeSecond = int(time.Hour.Seconds() * 24)
				}
				c.SetCookie(constants.TOKEN_COOKIE_NAME, tokenString, timeSecond, "/", "", false, true)
				models.SetKey(tokenString, 1)
				c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_OK, "", nil))
				return
			}
		}
		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_USERNAME_OR_PASSWORD_ERROR, constants.TIPS_USERNAME_OR_PASSWORD_ERROR, nil))
		return
	}
	c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_LOGIN_PARAM_NOT_VALID, constants.TIPS_LOGIN_PARAM_NOT_VALID, nil))
}

// 登出 能进入这里，说明token验证成功
func Logout(c *gin.Context) {
	tokenString, ok := c.Get("token")
	if ok {
		models.ClearKey(tokenString.(string))
		c.SetCookie(constants.TOKEN_COOKIE_NAME, "", -1, "/", "", false, false)
		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_OK, "", nil))
	} else {
		c.JSON(http.StatusInternalServerError, utils.ReturnJSON(constants.CODE_UNHANDLED_ERROR, "", nil))
	}
}

// 验证码
func VerifyCode(c *gin.Context) {

}
