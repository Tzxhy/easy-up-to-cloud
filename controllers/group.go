package controllers

import (
	"errors"
	"log"
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

func canReadGroupResource(uid, gid string) bool {
	commonGroup := models.GetCommonGroupResource()
	isCommon := utils.HasByFunc(commonGroup, func(m models.ResourceGroupItem) bool {
		return m.Gid == gid
	})
	if isCommon {
		return true
	}
	allGroup := models.GetAllGroupResource()
	item := utils.Find(allGroup, func(m models.ResourceGroupItem) bool {
		return m.Gid == gid
	})
	if item == nil {
		log.Fatal("item is nil: ", gid)
		return false
	}
	isAdmin := utils.HasByFunc(&models.AdminAccount, func(item models.AdminItem) bool {
		return item.Uid == uid
	})
	return isAdmin || utils.Has(&item.UserIds, uid)
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
	DirId string `json:"dir_id" form:"dir_id"`
}

func checkGroupIsVisibleToUser(uid, gid string) bool {
	// 检查该group可见性
	return models.GetGroupByIdAndUid(gid, uid) != nil
}

type ResourceGroupDirItemWithOp struct {
	models.ResourceGroupDirItem
	CanOper bool `json:"can_oper"`
}

// 获取目录信息
func GetGroupDir(c *gin.Context) {
	var getGroupDirReq GetGroupDirReq
	if c.ShouldBind(&getGroupDirReq) != nil {
		utils.ReturnParamNotValid(c)
		return
	}
	uid, _ := c.Get("uid")
	uidStr := uid.(string)

	if getGroupDirReq.DirId != "" {
		originResource := models.GetGroupResourceById(getGroupDirReq.DirId)
		if originResource == nil {
			c.JSON(http.StatusOK, utils.ReturnJSON(
				constants.CODE_GROUP_GET_DIR_WITH_ERROR_PARENT_TIPS.Code,
				constants.CODE_GROUP_GET_DIR_WITH_ERROR_PARENT_TIPS.Tip,
				nil,
			))
			return
		}

		if originResource.RType != models.GROUP_RESOURCE_DIR {
			c.JSON(http.StatusOK, utils.ReturnJSON(
				constants.CODE_GROUP_GET_DIR_WITH_ERROR_PARENT_FILE_TIPS.Code,
				constants.CODE_GROUP_GET_DIR_WITH_ERROR_PARENT_FILE_TIPS.Tip,
				nil,
			))
			return
		}
	}

	list := models.GetGroupDir(getGroupDirReq.DirId)
	commonGroup := models.GetCommonGroupResource()
	allGroup := models.GetAllGroupResource()
	var newList = make([]ResourceGroupDirItemWithOp, 0)
	currentIsAdmin := utils.Has(&models.AdminAccount, uidStr)
	for _, item := range *list {
		isCommon := false
		isOwnerUser := false
		if !currentIsAdmin {
			isCommon = utils.HasByFunc(*commonGroup, func(m models.ResourceGroupItem) bool {
				return m.Gid == item.Gid
			})
			gidItem := utils.Find(allGroup, func(m models.ResourceGroupItem) bool {
				return m.Gid == item.Gid
			})
			if gidItem == nil {
				log.Fatal(gidItem)
			}
			isOwnerUser = utils.Has(&gidItem.UserIds, uidStr)
		}
		if currentIsAdmin || isCommon || isOwnerUser {
			opItem := &ResourceGroupDirItemWithOp{
				item,
				currentIsAdmin || item.Uid == uidStr,
			}
			newList = append(newList, *opItem)
		}
	}
	c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_OK, "", &gin.H{
		"list": newList,
	}))
}

// type GroupCreateReq struct {
// 	Name string `json:"name" form:"name" binding:"required"`
// }

// 创建资源组
func groupCreate(c *gin.Context, name string) {
	// 只有管理员可以创建
	if !checkIsAdmin(c, true) {
		return
	}

	gid, err := models.CreateGroup(name, models.GroupTypeVisibleByUid)
	if err != nil {
		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_GROUP_CREATE_GROUP_HAS_TIPS.Code, constants.CODE_GROUP_CREATE_GROUP_HAS_TIPS.Tip, nil))
		return
	}
	c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_OK_TIPS.Code, "", &gin.H{
		"gid": gid,
	}))
}

type GroupCreateDirReq struct {
	ParentDid string `json:"parent_did" form:"parent_did"`
	Name      string `json:"name" form:"name" binding:"required"`
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

	if groupCreateDirReq.ParentDid == "" { // 创建分组
		groupCreate(c, groupCreateDirReq.Name)
		return
	}

	pid := groupCreateDirReq.ParentDid
	originResource := models.GetGroupResourceById(groupCreateDirReq.ParentDid)
	if originResource == nil {
		c.JSON(http.StatusOK, utils.ReturnJSON(
			constants.CODE_GROUP_CREATE_DIR_WITH_ERROR_PARENT_TIPS.Code,
			constants.CODE_GROUP_CREATE_DIR_WITH_ERROR_PARENT_TIPS.Tip,
			nil,
		))
		return
	}
	if originResource.RType != models.GROUP_RESOURCE_DIR {
		c.JSON(http.StatusOK, utils.ReturnJSON(
			constants.CODE_GROUP_CREATE_DIR_WITH_ERROR_PARENT_FILE_TIPS.Code,
			constants.CODE_GROUP_CREATE_DIR_WITH_ERROR_PARENT_FILE_TIPS.Tip,
			nil,
		))
		return
	}

	rid, err := models.CreateGroupDir(originResource.Gid, pid, groupCreateDirReq.Name, uid.(string))
	if err != nil {
		log.Print("GroupCreateDir model err: ", err)
		c.JSON(http.StatusOK, utils.ReturnJSON(
			constants.CODE_GROUP_CREATE_DIR_HAS_NAME_TIPS.Code,
			constants.CODE_GROUP_CREATE_DIR_HAS_NAME_TIPS.Tip,
			nil,
		))
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

	successNum := 0
	for _, group := range *groups {
		var isSuccess = false
		if utils.Has(&setGroupAccountReq.Groups, group.Gid) { // 有，则添加
			isSuccess, _ = models.SetUidResourceGroup(group.Gid, *utils.Unique(
				append(group.UserIds, setGroupAccountReq.Uid),
			))

		} else { // 无，则去除
			uids := *utils.Filter(&group.UserIds, func(t string) bool {
				return t != setGroupAccountReq.Uid
			})
			isSuccess, _ = models.SetUidResourceGroup(group.Gid, *utils.Unique(
				uids,
			))
		}
		if isSuccess {
			successNum = successNum + 1
		}
	}

	c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_OK_TIPS.Code, "", &gin.H{
		"successNum": successNum,
	}))
}

type ShareToGroupReq struct {
	// Fid
	Fid        string `json:"fid" form:"fid" binding:"required"`
	Name       string `json:"name" form:"name" binding:"required"`
	ParentDid  string `json:"parent_did" form:"parent_did" binding:"required"`
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

	resourceItem := models.GetGroupResourceById(shareToGroupReq.ParentDid)
	if resourceItem == nil {
		c.JSON(http.StatusOK, utils.ReturnJSON(
			constants.CODE_GROUP_RESOURCE_NOT_FOUND_TIPS.Code,
			constants.CODE_GROUP_RESOURCE_NOT_FOUND_TIPS.Tip,
			nil,
		))
		return
	}

	canRead := canReadGroupResource(uid.(string), resourceItem.Gid)
	if !canRead {
		c.JSON(http.StatusOK, utils.ReturnJSON(
			constants.CODE_GROUP_NOT_FOUND_TIPS.Code,
			constants.CODE_GROUP_NOT_FOUND_TIPS.Tip,
			nil,
		))
		return
	}

	rid, err := models.ShareFileToGroup(
		resourceItem.Gid,
		shareToGroupReq.Fid,
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

const (
	GroupResourceDelete = "delete"
	GroupResourceRename = "rename"
	GroupResourceExpire = "expire"
)

type OperationGroupResourceReq struct {
	// 资源id
	Rid string `json:"rid" form:"rid" binding:"required"`
	// 操作
	Oper string `json:"oper" form:"oper" binding:"required"`
	// rename时新名称
	NewName string `json:"new_name" form:"new_name"`
	// move 时新父目录
	ParentDid string `json:"parent_did" form:"parent_did"`
	// 修改过期日期时
	ExpireDate string `json:"expire_date" form:"expire_date"`
}
type OperationGroupResourceReqInner struct {
	OperationGroupResourceReq
	Uid string
}

// 操作已经分享到资源组的资源
func OperationGroupResource(c *gin.Context) {
	var operationGroupResourceReq OperationGroupResourceReq
	if c.ShouldBind(&operationGroupResourceReq) != nil {
		utils.ReturnParamNotValid(c)
		return
	}
	handleOperationGroupResourceDispatch.name = operationGroupResourceReq.Oper
	uid, _ := c.Get("uid")
	operationGroupResourceReqInner := &OperationGroupResourceReqInner{
		operationGroupResourceReq,
		uid.(string),
	}
	err := handleOperationGroupResourceDispatch.handle(operationGroupResourceReqInner)
	if err == nil {
		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_OK, "", nil))
		return
	}
	c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_GROUP_OPERATION_FAILED_TIPS.Code, err.Error(), nil))
}

type HandleOperationGroupResourceDispatch struct {
	name       string
	strategies []HandleOperationGroupResourceStrategy
}
type HandleOperationGroupResourceStrategy struct {
	canHandle func(string) bool
	handle    func(*OperationGroupResourceReqInner) error
}

func (h *HandleOperationGroupResourceDispatch) handle(o *OperationGroupResourceReqInner) error {
	var err error
	for _, strategy := range h.strategies {
		if strategy.canHandle(o.Oper) {
			err = strategy.handle(o)
			break
		}
	}
	return err
}
func HandleDeleteRid(o *OperationGroupResourceReqInner) error {
	// 定义为只有admin用户可以删除文件夹
	// user用户可以删除自己分享的文件
	isAdmin := utils.Has(&models.AdminAccount, o.Uid)
	resourceItem := models.GetGroupResourceById(o.Rid)

	if resourceItem.RType == models.GROUP_RESOURCE_FILE { // 文件，直接删
		if isAdmin {
			models.DeleteResourceByRid(o.Rid)
		}
		return nil
	} else { // 文件夹，则递归删
		list := models.GetGroupDir(o.Rid)

		oList := utils.Map(list, func(resource models.ResourceGroupDirItem) *OperationGroupResourceReqInner {
			return &OperationGroupResourceReqInner{
				OperationGroupResourceReq{
					Rid: resource.Rid,
				},
				o.Uid,
			}
		})
		for _, child := range *oList {
			HandleDeleteRid(child)
		}

		if isAdmin {
			models.DeleteResourceByRid(o.Rid)
		}

	}
	return nil
}

var handleOperationGroupResourceDispatch = &HandleOperationGroupResourceDispatch{
	name: "",
	strategies: []HandleOperationGroupResourceStrategy{
		// 删除
		{
			canHandle: func(name string) bool {
				return name == GroupResourceDelete
			},
			handle: HandleDeleteRid,
		},
		// 重命名
		{
			canHandle: func(name string) bool {
				return name == GroupResourceRename
			},
			handle: func(o *OperationGroupResourceReqInner) error {
				// 检查
				if o.NewName == "" {
					return errors.New("名称不能为空")
				}
				_, err := models.RenameResourceByUidAndRid(o.Uid, o.Rid, o.NewName)
				return err
			},
		},
		// 有效期
		{
			canHandle: func(name string) bool {
				return name == GroupResourceExpire
			},
			handle: func(o *OperationGroupResourceReqInner) error {
				_, err := models.ExpireChangeResourceByUidAndRid(o.ExpireDate, o.Uid, o.Rid)
				return err
			},
		},
	},
}

type SearchGroupResourceReq struct {
	Name string `json:"name" form:"name" binding:"required"`
}

// 搜索
func SearchGroupResource(c *gin.Context) {
	var searchGroupResourceReq SearchGroupResourceReq
	if c.ShouldBindQuery(&searchGroupResourceReq) != nil {
		utils.ReturnParamNotValid(c)
		return
	}
	uid, _ := c.Get("uid")
	uidStr := uid.(string)

	list := models.SearchResourceByName(searchGroupResourceReq.Name)

	var newList = make([]ResourceGroupDirItemWithOp, 0)
	commonGroup := models.GetCommonGroupResource()
	currentIsAdmin := utils.Has(&models.AdminAccount, uidStr)
	allGroup := models.GetAllGroupResource()
	for _, item := range *list {
		isCommon := false
		isOwnerUser := false
		if !currentIsAdmin {
			isCommon = utils.HasByFunc(*commonGroup, func(m models.ResourceGroupItem) bool {
				return m.Gid == item.Gid
			})
			gidItem := utils.Find(allGroup, func(m models.ResourceGroupItem) bool {
				return m.Gid == item.Gid
			})
			if gidItem == nil {
				log.Fatal(gidItem)
			}
			isOwnerUser = utils.Has(&gidItem.UserIds, uidStr)
		}

		if currentIsAdmin || isCommon || isOwnerUser {
			opItem := &ResourceGroupDirItemWithOp{
				item,
				item.Uid == uidStr,
			}
			newList = append(newList, *opItem)
		}
	}
	c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_OK, "", &gin.H{
		"list": newList,
	}))

}

// 预览
func PreviewGroupResource(c *gin.Context) {
	var downloadGroupResourceReq DownloadGroupResourceReq
	if c.ShouldBindQuery(&downloadGroupResourceReq) != nil {
		utils.ReturnParamNotValid(c)
		return
	}
	resourceItem := models.GetGroupResourceById(downloadGroupResourceReq.Rid)
	if resourceItem == nil {
		c.JSON(http.StatusOK, utils.ReturnJSON(
			constants.CODE_GROUP_RESOURCE_NOT_FOUND_TIPS.Code,
			constants.CODE_GROUP_RESOURCE_NOT_FOUND_TIPS.Tip,
			nil,
		))
		return
	}

	uid, _ := c.Get("uid")
	currentIsAdmin := utils.Has(&models.AdminAccount, uid.(string))

	if !currentIsAdmin && !canReadGroupResource(uid.(string), resourceItem.Gid) {
		c.JSON(http.StatusOK, utils.ReturnJSON(
			constants.CODE_GROUP_PREVIEW_FILE_NO_PERMISSION_TIPS.Code,
			constants.CODE_GROUP_PREVIEW_FILE_NO_PERMISSION_TIPS.Tip,
			nil,
		))
		return
	}
	file := models.GetFileById(resourceItem.Fid)
	if file == nil {
		c.JSON(http.StatusOK, utils.ReturnJSON(
			constants.CODE_GROUP_RESOURCE_NOT_FOUND_TIPS.Code,
			constants.CODE_GROUP_RESOURCE_NOT_FOUND_TIPS.Tip,
			nil,
		))
		return
	}
	previewFile(resourceItem.Fid, file.OwnerId, c)
}

type DownloadGroupResourceReq struct {
	Rid string `form:"rid" json:"rid" binding:"required"`
}

// 下载
func DownloadGroupResource(c *gin.Context) {
	var downloadGroupResourceReq DownloadGroupResourceReq
	if c.ShouldBindQuery(&downloadGroupResourceReq) != nil {
		utils.ReturnParamNotValid(c)
		return
	}
	resourceItem := models.GetGroupResourceById(downloadGroupResourceReq.Rid)
	if resourceItem == nil {
		c.JSON(http.StatusOK, utils.ReturnJSON(
			constants.CODE_GROUP_RESOURCE_NOT_FOUND_TIPS.Code,
			constants.CODE_GROUP_RESOURCE_NOT_FOUND_TIPS.Tip,
			nil,
		))
		return
	}

	uid, _ := c.Get("uid")
	currentIsAdmin := utils.Has(&models.AdminAccount, uid.(string))
	log.Print("models.AdminAccount: ", models.AdminAccount)
	log.Print("uid: ", uid)

	if !currentIsAdmin && !canReadGroupResource(uid.(string), resourceItem.Gid) {
		c.JSON(http.StatusOK, utils.ReturnJSON(
			constants.CODE_GROUP_PREVIEW_FILE_NO_PERMISSION_TIPS.Code,
			constants.CODE_GROUP_PREVIEW_FILE_NO_PERMISSION_TIPS.Tip,
			nil,
		))
		return
	}
	file := models.GetFileById(resourceItem.Fid)
	if file == nil {
		c.JSON(http.StatusOK, utils.ReturnJSON(
			constants.CODE_GROUP_RESOURCE_NOT_FOUND_TIPS.Code,
			constants.CODE_GROUP_RESOURCE_NOT_FOUND_TIPS.Tip,
			nil,
		))
		return
	}
	downloadFile(resourceItem.Fid, file.OwnerId, c)
}

func GroupUserConfig(c *gin.Context) {

}

type MoveResourcesReq struct {
	Rids      []string `json:"rids" form:"rids" binding:"required"`
	ParentDid string   `json:"parent_did" form:"parent_did" binding:"required"`
}

func MoveResources(c *gin.Context) {
	var moveResourcesReq MoveResourcesReq
	if c.ShouldBind(&moveResourcesReq) != nil {
		utils.ReturnParamNotValid(c)
		return
	}

	uid, _ := c.Get("uid")
	uidStr := uid.(string)
	isAdmin := utils.Has(&models.AdminAccount, uidStr)
	if isAdmin { // 管理员直接全权
		count := models.MoveMultiGroups(moveResourcesReq.Rids, moveResourcesReq.ParentDid)

		c.JSON(http.StatusOK, utils.ReturnJSON(
			constants.CODE_OK,
			"",
			&gin.H{
				"successNum": count,
			},
		))
		return
	}
}
