package controllers

import (
	"log"
	"net/http"
	"path/filepath"

	"gitee.com/tzxhy/web/constants"
	"gitee.com/tzxhy/web/utils"
	"github.com/gin-gonic/gin"
)

func ListDir(c *gin.Context) {

}

// 上传文件
func UploadFile(c *gin.Context) {
	myFile, _ := c.FormFile("file")
	filePath := filepath.Join(constants.UPLOAD_PATH, myFile.Filename)
	err := c.SaveUploadedFile(myFile, filePath)
	if err != nil {
		log.Print(err)
	}
	log.Print(filePath)
	c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_OK, "", nil))
}

// 上传文件
func DownloadFile(c *gin.Context) {

}

// 预览
func PreviewFile(c *gin.Context) {

}

// 重命名
func RenameFile(c *gin.Context) {

}

func DeleteFile(c *gin.Context) {

}

// 分享资源
func ShareFileToGroup(c *gin.Context) {

}
