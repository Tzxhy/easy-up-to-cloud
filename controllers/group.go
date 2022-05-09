package controllers

import (
	"net/http"

	"gitee.com/tzxhy/web/constants"
	"gitee.com/tzxhy/web/models"
	"gitee.com/tzxhy/web/utils"
	"github.com/gin-gonic/gin"
)

// type GetMyGroupsReq struct {}
// 获取我的资源组列表
func GetMyGroups(c *gin.Context) {
	uid, _ := c.Get("uid")

	groupItems := models.GetResourceGroup(uid.(string))
	c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_OK, "", &gin.H{
		"groups": groupItems,
	}))

}
