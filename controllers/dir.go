package controllers

import (
	"fmt"
	"net/http"

	"gitee.com/tzxhy/web/constants"
	"gitee.com/tzxhy/web/models"
	"gitee.com/tzxhy/web/utils"
	"github.com/gin-gonic/gin"
)

type NewDirInfo struct {
	ParentDirId *int   `json:"parent_id"`
	Name        string `json:"name" binding:"required"`
}

// 创建目录
func CreateDir(c *gin.Context) {
	var newDirInfo NewDirInfo
	err := c.ShouldBindJSON(&newDirInfo)
	uid, _ := c.Get("uid")
	parnetDidNum := 0
	
	if newDirInfo.ParentDirId != nil {
		parnetDidNum = *newDirInfo.ParentDirId
	} else {
		parnetDidNum = -1
	}
	did, err := models.AddDir(uid.(int), newDirInfo.Name, parnetDidNum)
	if err == nil {
		fmt.Print(did)
		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_OK, "", &gin.H{
			"did": did,
		}))
	} else {
		
		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_UNHANDLED_ERROR, err.Error(), nil))
	}
}

// 获取目录信息
func GetDir(c *gin.Context) {

}

// 查找目录或者文件
func SearchFileOrDir(c *gin.Context) {

}
