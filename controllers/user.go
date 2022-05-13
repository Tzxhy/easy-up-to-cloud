package controllers

import (
	"net/http"

	"gitee.com/tzxhy/web/constants"
	"gitee.com/tzxhy/web/models"
	"gitee.com/tzxhy/web/utils"
	"github.com/gin-gonic/gin"
)

type UserInfo struct {
	Uid      string `json:"uid"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
}

func GetUserInfo(c *gin.Context) {
	uid, _ := c.Get("uid")
	userModel := models.GetUserById(uid.(string))
	isAdmin := utils.Has(models.GetAdminAccount(), uid.(string))
	userInfo := &UserInfo{
		Uid:      userModel.Uid,
		Username: userModel.Username,
		IsAdmin:  isAdmin,
	}
	c.JSON(http.StatusOK, utils.ReturnJSON(
		constants.CODE_OK,
		"",
		&gin.H{
			"user_info": userInfo,
		},
	))
}
