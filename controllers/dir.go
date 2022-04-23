package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type NewDirInfo struct {
	ParentDirId int    `form:"parent_id" binding:"required"`
	Name        string `form:"name" binding:"required"`
}

// 创建目录
func CreateDir(c *gin.Context) {
	var newDirInfo NewDirInfo
	err := c.ShouldBind(&newDirInfo)
	if err == nil {
		fmt.Print(newDirInfo)
	} else {
		fmt.Print(err)
	}
}

// 获取目录信息
func GetDir(c *gin.Context) {

}

// 查找目录或者文件
func SearchFileOrDir(c *gin.Context) {

}
