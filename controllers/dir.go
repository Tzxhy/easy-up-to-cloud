package controllers

import (
	"net/http"

	"gitee.com/tzxhy/web/constants"
	"gitee.com/tzxhy/web/models"
	"gitee.com/tzxhy/web/utils"
	"github.com/gin-gonic/gin"
)

type NewDirInfo struct {
	ParentDirId string `json:"parent_id"`
	Name        string `json:"name" binding:"required"`
}

// 创建目录
func CreateDir(c *gin.Context) {
	var newDirInfo NewDirInfo
	err := c.ShouldBind(&newDirInfo)
	if err == nil {
		uid, _ := c.Get("uid")
		did, err := models.AddDir(uid.(string), newDirInfo.Name, newDirInfo.ParentDirId)
		if err == nil {
			c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_OK, "", &gin.H{
				"did": did,
			}))
		} else {
			c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_UNHANDLED_ERROR, err.Error(), nil))
		}
	} else {
		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_CREATE_DIR_PARAM_NOT_VALID, constants.TIPS_CREATE_DIR_PARAM_NOT_VALID, nil))
	}
}

type GetDirInfo struct {
	DirId string `form:"dir_id" binding:"required"`
}

// 获取目录信息
func GetDir(c *gin.Context) {
	var getDirInfo GetDirInfo
	if err := c.ShouldBindQuery(&getDirInfo); err == nil {
		uid, _ := c.Get("uid")
		dirInfo := models.GetDir(getDirInfo.DirId, uid.(string))
		if dirInfo != nil {
			c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_OK, "", &gin.H{
				"did":         dirInfo.Did,
				"dirname":     dirInfo.Dirname,
				"parent_did":  dirInfo.ParentDiD,
				"create_date": dirInfo.CreateDate,
			}))
		} else {
			c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_QUERY_DIR_INFO_WITH_EMPTY_RES, constants.TIPS_QUERY_DIR_INFO_WITH_EMPTY_RES, nil))
		}
	} else {
		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_PARAMS_NOT_VALID, constants.TIPS_COMMON_PARAM_NOT_VALID, nil))
	}
}

// 查找目录或者文件
func SearchFileOrDir(c *gin.Context) {

}
