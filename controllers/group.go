package controllers

import (
	"net/http"

	"gitee.com/tzxhy/web/constants"
	"gitee.com/tzxhy/web/models"
	"gitee.com/tzxhy/web/utils"
	"github.com/gin-gonic/gin"
)

const IS_ADMIN_PREFIX = "is_admin"

func isAdminAccount(uid string) bool {
	isAdmin := models.GetKey(IS_ADMIN_PREFIX + uid)
	isAd := false
	if isAdmin == nil { // 未保存过
		admins := models.GetAdminUser()
		if utils.Some(*admins, func(t models.User) bool {
			return t.Uid == uid
		}) { // 有
			models.SetKey(IS_ADMIN_PREFIX+uid, true)
			return true
		} else {
			models.SetKey(IS_ADMIN_PREFIX+uid, false)
			return false
		}
	} else {
		isAd = isAdmin.(bool)
		return isAd
	}
}

// 获取我的资源组列表
func GetMyGroups(c *gin.Context) {
	uid, _ := c.Get("uid")
	// 判断是否是admin账户，是的话，可以查看所有
	isAdmin := isAdminAccount(uid.(string))
	var groupItems *[]models.ResourceGroupItem
	if isAdmin {
		groupItems = models.GetAllResourceGroup()
	} else {
		groupItems = models.GetResourceGroup(uid.(string))
	}

	c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_OK, "", &gin.H{
		"groups": groupItems,
	}))
}

type GroupCreateDirReq struct {
	// group_id
	Gid       string  `json:"gid" form:"gid" binding:"required"`
	ParentDid *string `json:"parent_did" form:"parent_did"`
	Name      string  `json:"name" form:"name" binding:"required"`
}

// 创建文件夹
func GroupCreateDir(c *gin.Context) {
	var groupCreateDirReq GroupCreateDirReq
	if c.ShouldBind(&groupCreateDirReq) != nil {
		utils.ReturnParamNotValid(c)
		return
	}

	// 不存在的gid,pid，则错误
}

func GroupEmpty(c *gin.Context) {

}
