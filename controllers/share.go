package controllers

import (
	"net/http"
	"time"

	"gitee.com/tzxhy/web/constants"
	"gitee.com/tzxhy/web/models"
	"gitee.com/tzxhy/web/utils"
	"github.com/gin-gonic/gin"
)

type CreateShareReq struct {
	Fid        string `json:"fid" form:"fid"`
	Did        string `json:"did" form:"did"`
	Name       string `json:"name" form:"name" binding:"required"`
	Password   string `json:"password" form:"password"`
	ExpireDate int64  `json:"expire_date" form:"expire_date"`
}

// 创建分享
func CreateShare(c *gin.Context) {
	var createShareReq CreateShareReq
	if c.ShouldBind(&createShareReq) != nil {
		utils.ReturnParamNotValid(c)
		return
	}
	expireDate := int64(-1)
	if createShareReq.ExpireDate > 0 {
		expireDate = createShareReq.ExpireDate
	}
	uid, _ := c.Get("uid")
	_, err := models.AddShareItem(
		createShareReq.Fid,
		createShareReq.Did,
		createShareReq.Name,
		uid.(string),
		createShareReq.Password,
		expireDate,
	)
	if err == nil {
		c.JSON(http.StatusOK, utils.ReturnJSON(
			constants.CODE_OK,
			"",
			nil,
		))
		return
	}
	c.JSON(http.StatusOK, utils.ReturnJSON(
		constants.CODE_SHARE_FAIL_TIPS.Code,
		constants.CODE_SHARE_FAIL_TIPS.Tip+": "+err.Error(),
		nil,
	))
}

type ShareItemDetail struct {
	models.ShareItem
	ShareUserName string `json:"share_user"`
	CanOper       bool   `json:"can_oper"`
	NeedPassword  bool   `json:"lock"`
}

// 查看分享列表
func GetShareList(c *gin.Context) {
	uid, hasUid := c.Get("uid")
	list := models.GetAllShareItems()
	if list != nil {
		users := models.GetUserByIds(
			*utils.Unique(*utils.Map(list, func(item models.ShareItem) string {
				return item.UserId
			})),
		)
		shareItemDetails := utils.Map(list, func(item models.ShareItem) ShareItemDetail {
			return ShareItemDetail{
				item,
				utils.Find(users, func(user models.User) bool {
					return user.Uid == item.UserId
				}).Username,
				hasUid && item.UserId == uid.(string),
				item.Password != "",
			}
		})
		c.JSON(http.StatusOK, utils.ReturnJSON(
			constants.CODE_OK,
			"",
			&gin.H{
				"list": shareItemDetails,
			},
		))
	}
}

type GetShareDetailReq struct {
	Sid      string `json:"sid" form:"sid" binding:"required"`
	Password string `json:"password" form:"password"`
}

type ShareDetailItem struct {
	models.File
	ShareUserName string `json:"share_user"`
	Url           string `json:"url"`
}

// 查看分享详情
// 已登录用户，判断是否是分享本人；
// 未登录，判断是否有该item的cookie（从请求带password来）
func GetShareDetail(c *gin.Context) {
	var getShareDetailReq GetShareDetailReq
	if c.ShouldBindQuery(&getShareDetailReq) != nil {
		utils.ReturnParamNotValid(c)
		return
	}

	item := models.GetShareItem(getShareDetailReq.Sid)

	if item == nil {
		c.JSON(http.StatusOK, utils.ReturnJSON(
			constants.CODE_SHARE_ITEM_NOT_FOUND_TIPS.Code,
			constants.CODE_SHARE_ITEM_NOT_FOUND_TIPS.Tip,
			nil,
		))
		return
	}

	uid, hasUid := c.Get("uid")

	if hasUid { // 已登录
		isSharer := uid.(string) == item.UserId
		if !isSharer && getShareDetailReq.Password != item.Password {
			c.JSON(http.StatusOK, utils.ReturnJSON(
				constants.CODE_SHARE_ITEM_NEED_PASSWORD_TIPS.Code,
				constants.CODE_SHARE_ITEM_NEED_PASSWORD_TIPS.Tip,
				nil,
			))
			return
		}
	} else { // 未登录
		if getShareDetailReq.Password == item.Password {
			tokenString, _ := utils.GenShareItemToken(item.Sid)
			timeSecond := int(time.Hour.Seconds() * 24)
			// 添加 cookie 以便后续资源访问时可以接续
			c.SetCookie(constants.SHARE_ITEM_PREFIX_COOKIE_NAME+item.Sid, tokenString, timeSecond, "/", "", false, true)
		} else {
			c.JSON(http.StatusOK, utils.ReturnJSON(
				constants.CODE_SHARE_ITEM_NEED_PASSWORD_TIPS.Code,
				constants.CODE_SHARE_ITEM_NEED_PASSWORD_TIPS.Tip,
				nil,
			))
			return
		}
	}

	if item.Fid != "" { // 文件分享
		file := models.GetFileById(item.Fid)
		user := models.GetUserById(file.UserId)
		var ret = ShareDetailItem{
			*file,
			user.Username,
			getDownloadUrl(item.Sid),
		}
		c.JSON(http.StatusOK, utils.ReturnJSON(
			constants.CODE_OK,
			"",
			&gin.H{
				"type": 2,
				"file": ret,
			},
		))
		return
	} else if item.Did != "" {
		list := models.GetFileList(item.Did, item.UserId)
		users := models.GetUserByIds(
			*utils.Unique(*utils.Map(list, func(item models.File) string {
				return item.UserId
			})),
		)
		newList := utils.Map(list, func(itemIn models.File) ShareDetailItem {
			url := getDownloadUrl(item.Sid) + "&fid=" + itemIn.Fid

			return ShareDetailItem{
				itemIn,
				utils.Find(users, func(user models.User) bool {
					return user.Uid == itemIn.UserId
				}).Username,
				url,
			}
		})
		c.JSON(http.StatusOK, utils.ReturnJSON(
			constants.CODE_OK,
			"",
			&gin.H{
				"type": 1,
				"dir":  newList,
			},
		))
	}
}

func getDownloadUrl(sid string) string {
	return "/api/v1/share/download?sid=" + sid
}

type ShareDownloadReq struct {
	Sid      string `json:"sid" form:"sid" binding:"required"`
	Fid      string `json:"fid" form:"fid"`
	Password string `json:"password" form:"password"`
}

func ShareDownload(c *gin.Context) {
	var shareDownloadReq ShareDownloadReq
	if c.ShouldBindQuery(&shareDownloadReq) != nil {
		utils.ReturnParamNotValid(c)
		return
	}

	item := models.GetShareItem(shareDownloadReq.Sid)
	if item == nil {
		c.JSON(http.StatusOK, utils.ReturnJSON(
			-1,
			"没找到",
			nil,
		))
		return
	}
	if item.Fid != "" { // 单文件
		downloadFile(item.Fid, item.UserId, item.Name, c)
	} else { // 目录形式分享
		if shareDownloadReq.Fid == "" {
			c.JSON(http.StatusOK, utils.ReturnJSON(
				-1,
				"参数不对",
				nil,
			))
			return
		}
		file := models.GetFileById(shareDownloadReq.Fid)
		if file == nil {
			c.JSON(http.StatusOK, utils.ReturnJSON(
				-1,
				"没找到文件",
				nil,
			))
			return
		}
		downloadFile(file.Fid, item.UserId, file.Filename, c)
	}
}

func DeleteShare(c *gin.Context) {
	var deleteReq GetShareDetailReq
	if c.ShouldBind(&deleteReq) != nil {
		utils.ReturnParamNotValid(c)
		return
	}
	uid, _ := c.Get("uid")
	ok := models.DeleteShare(deleteReq.Sid, uid.(string))
	c.JSON(http.StatusOK, utils.ReturnJSON(
		constants.CODE_OK,
		"",
		&gin.H{
			"success": ok,
		},
	))
}
