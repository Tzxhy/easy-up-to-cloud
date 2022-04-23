package controllers

import (
	"log"

	"gitee.com/tzxhy/web/constants"
	"gitee.com/tzxhy/web/models"
	"gitee.com/tzxhy/web/utils"
	"github.com/gin-gonic/gin"
)

type LoginInfo struct {
	Username string `form:"username"`
	Password string `form:"Password"`
}

// 登录
func Login(c *gin.Context) {
	var loginInfo LoginInfo
	if c.ShouldBind(&loginInfo) == nil {
		if loginInfo.Username != "" && loginInfo.Password != "" {
			userInfo := models.GetUserByNameAndPassword(loginInfo.Username, loginInfo.Password)
			log.Print(userInfo)
			if userInfo.Uid != 0 {

				tokenString, _ := utils.GenToken(userInfo.Username)
				c.SetCookie(constants.TOKEN_COOKIE_NAME, tokenString, 60000, "/", "", false, false)
				log.Print(tokenString)
				models.SetKey(tokenString, 1)
				log.Print("登录成功")
				c.String(200, "Success")
			}
		}
	}
	c.String(401, "")

}

type LogoutPerson struct {
	Username string `form:"username"`
}
type LogoutInfo struct {
	Token string `form:"token"`
}

// 登出
func Logout(c *gin.Context) {
	var logoutInfo LogoutInfo
	if c.ShouldBind(&logoutInfo) == nil {
		log.Println(logoutInfo.Token)
		if models.GetKey(logoutInfo.Token) != nil {

			log.Println("清除Token: ", logoutInfo.Token)
		}
	}

	c.String(200, "Success")
}

// 验证码
func VerifyCode(c *gin.Context) {

}
