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

func checkIsAdmin(c *gin.Context, needRejectWhenNot bool) bool {
	uid, _ := c.Get("uid")
	isAdmin := isAdminAccount(uid.(string))
	if needRejectWhenNot && !isAdmin {
		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_GROUP_IS_NOT_ADMIN_TIPS.Code, constants.CODE_GROUP_IS_NOT_ADMIN_TIPS.Tip, nil))
	}
	return isAdmin
}

// 获取我的资源组列表
func GetMyGroups(c *gin.Context) {
	uid, _ := c.Get("uid")
	// 判断是否是admin账户，是的话，可以查看所有
	isAdmin := checkIsAdmin(c, false)
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
	uid, _ := c.Get("uid")
	if !checkIsAdmin(c, true) {
		return
	}

	var groupCreateDirReq GroupCreateDirReq
	if c.ShouldBind(&groupCreateDirReq) != nil {
		utils.ReturnParamNotValid(c)
		return
	}

	pid := ""
	if groupCreateDirReq.ParentDid != nil {
		pid = *groupCreateDirReq.ParentDid
	}
	// 不存在的gid,pid，则错误
	rid, err := models.CreateGroupDir(groupCreateDirReq.Gid, pid, groupCreateDirReq.Name, uid.(string))
	if err != nil {
		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_UNHANDLED_ERROR, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_OK_TIPS.Code, "", &gin.H{
		"rid": rid,
	}))
}

type SetGroupAccountReq struct {
	// 用户id
	Uid string `json:"uid" form:"uid" binding:"required"`
	// 所属资源组
	Groups []string `json:"groups" form:"groups"`
	// 是否设置为管理员账号
	Admin *bool `json:"is_admin" form:"is_admin"`
}

// 设置账户信息
func SetGroupAccount(c *gin.Context) {
	if !checkIsAdmin(c, true) {
		return
	}
	var setGroupAccountReq SetGroupAccountReq
	if c.ShouldBind(&setGroupAccountReq) != nil {
		utils.ReturnParamNotValid(c)
		return
	}
	// 设置管理员权限
	uid, _ := c.Get("uid")
	if setGroupAccountReq.Admin != nil { // 非空，则需要
		isAdmin := *setGroupAccountReq.Admin
		models.DeleteOrInsertAdminAccount(uid.(string), isAdmin)
	}
}

func GroupEmpty(c *gin.Context) {

}
