package controllers

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

type NewDirInfo struct {
	ParentDirId *int   `json:"parent_id" validate:"exists, numeric"`
	Name        string `json:"name" binding:"required"`
}

// 创建目录
func CreateDir(c *gin.Context) {
	var newDirInfo NewDirInfo
	err := c.ShouldBindJSON(&newDirInfo)
	log.Print(newDirInfo)
	log.Print(newDirInfo.ParentDirId)
	log.Print(*newDirInfo.ParentDirId)
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
