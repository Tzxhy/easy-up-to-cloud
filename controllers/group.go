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

// 检查是否是admin账户
func checkIsAdmin(c *gin.Context, needRejectWhenNot bool) bool {
	uid, _ := c.Get("uid")
	isAdmin := isAdminAccount(uid.(string))
	if needRejectWhenNot && !isAdmin {
		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_GROUP_IS_NOT_ADMIN_TIPS.Code, constants.CODE_GROUP_IS_NOT_ADMIN_TIPS.Tip, nil))
	}
	return isAdmin
}

// 获取资源组列表
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

type GetGroupDirReq struct {
	Gid   string  `json:"gid" form:"gid" binding:"required"`
	DirId *string `json:"dir_id" form:"dir_id"`
}

// 获取目录信息
func GetGroupDir(c *gin.Context) {
	var getGroupDirReq GetGroupDirReq
	if c.ShouldBind(&getGroupDirReq) != nil {
		utils.ReturnParamNotValid(c)
		return
	}

	list := models.GetGroupDir(getGroupDirReq.Gid, utils.GetStringOrEmptyFromPtr(getGroupDirReq.DirId))
	c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_OK, "", &gin.H{
		"list": list,
	}))
}

type GroupCreateReq struct {
	Name string `json:"name" form:"name" binding:"required"`
}

// 创建资源组
func GroupCreate(c *gin.Context) {
	if !checkIsAdmin(c, true) {
		return
	}

	var groupCreateReq GroupCreateReq
	if c.ShouldBind(&groupCreateReq) != nil {
		utils.ReturnParamNotValid(c)
		return
	}

	gid, err := models.CreateGroup(groupCreateReq.Name)
	if err != nil {
		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_GROUP_CREATE_GROUP_HAS_TIPS.Code, constants.CODE_GROUP_CREATE_GROUP_HAS_TIPS.Tip, nil))
		return
	}
	c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_OK_TIPS.Code, "", &gin.H{
		"gid": gid,
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

	pid := utils.GetStringOrEmptyFromPtr(groupCreateDirReq.ParentDid)

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
	// 设置账户资源组
	groups := models.GetAllResourceGroup()
	var filteredGroups []models.ResourceGroupItem
	for _, group := range *groups {
		if utils.Has(setGroupAccountReq.Groups, group.Gid) {
			filteredGroups = append(filteredGroups, group)
		}
	}
	succNum := 0
	for _, group := range filteredGroups {
		succ, _ := models.SetUidResourceGroup(group.Gid, *utils.Unique(
			append(group.UserIds, setGroupAccountReq.Uid),
		))
		if succ {
			succNum = succNum + 1
		}
	}
	c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_OK_TIPS.Code, "", &gin.H{
		"succNum": succNum,
	}))
}

type ShareToGroupReq struct {
	Gid string `json:"gid" form:"gid" binding:"required"`
	// Fid 与 Did 二选一
	Fid        string `json:"fid" form:"fid"`
	Did        string `json:"did" form:"did"`
	Name       string `json:"name" form:"name" binding:"required"`
	ParentDid  string `json:"parent_did" form:"parent_did"`
	RType      uint8  `json:"r_type" form:"r_type" binding:"required"`
	ExpireDate string `json:"expire_date" form:"expire_date"`
}

// 分享资源到资源组
func ShareToGroup(c *gin.Context) {
	var shareToGroupReq ShareToGroupReq
	if c.ShouldBind(&shareToGroupReq) != nil {
		utils.ReturnParamNotValid(c)
		return
	}

	uid, _ := c.Get("uid")

	rid, err := models.ShareDirOrFileToGroup(
		shareToGroupReq.Gid,
		shareToGroupReq.Fid,
		shareToGroupReq.Did,
		shareToGroupReq.Name,
		uid.(string),
		shareToGroupReq.ParentDid,
		shareToGroupReq.ExpireDate,
		shareToGroupReq.RType,
	)
	if err != nil {
		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_UNHANDLED_ERROR, err.Error(), nil))
	} else {
		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_OK, "", &gin.H{
			"rid": rid,
		}))
	}
}

func GroupEmpty(c *gin.Context) {

}
