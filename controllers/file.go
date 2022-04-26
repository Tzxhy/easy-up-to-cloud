package controllers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"gitee.com/tzxhy/web/constants"
	"gitee.com/tzxhy/web/models"
	"gitee.com/tzxhy/web/utils"
	"github.com/gin-gonic/gin"
)

func ListDir(c *gin.Context) {

}

type UploadFileReq struct {
	ParentDid *string `form:"parent_did"`
}

// 上传文件
func UploadFile(c *gin.Context) {
	myFile, err := c.FormFile("file")
	if err != nil {
		log.Print(err)
	}

	uid, _ := c.Get("uid")
	did := ""
	var uploadFileReq UploadFileReq
	c.Bind(&uploadFileReq)
	if uploadFileReq.ParentDid != nil {
		did = *uploadFileReq.ParentDid
	}
	// log.Print(myFile)
	// log.Print(uid)
	// log.Print(did)
	file := models.GetFileByName(myFile.Filename, uid.(string), did)
	if file != nil {
		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_FILENAME_HAS_BEEN_USED, constants.TIPS_FILENAME_HAS_BEEN_USED, nil))
		return
	}
	filePath := filepath.Join(constants.UPLOAD_PATH, uid.(string), myFile.Filename)
	utils.MakeSurePathExists(filepath.Dir(filePath))
	log.Print(filePath)
	err = c.SaveUploadedFile(myFile, filePath)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, utils.ReturnJSON(constants.CODE_UNHANDLED_ERROR, err.Error(), nil))
		return
	}
	fid, err := models.AddFile(
		uid.(string),
		did,
		myFile.Filename,
		uint64(myFile.Size),
		filePath,
	)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, utils.ReturnJSON(constants.CODE_UNHANDLED_ERROR, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_OK, "", &gin.H{
		"fid": fid,
	}))
}

type FileIdReq struct {
	Fid string `form:"fid" binding:"required"`
}

// 下载文件
func DownloadFile(c *gin.Context) {
	var FileIdReq FileIdReq
	if c.ShouldBindQuery(&FileIdReq) == nil {
		uid, _ := c.Get("uid")
		fileInfo := models.GetFile(FileIdReq.Fid, uid.(string))
		if fileInfo == nil {
			c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_CREATE_DIR_PARAM_NOT_VALID, constants.TIPS_FILE_NOT_EXIST, nil))
			return
		}
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", "attachment; filename="+fileInfo.Filename)
		c.Header("Content-Transfer-Encoding", "binary")
		c.File(fileInfo.FileRealPath)
		return
	}
	c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_PARAMS_NOT_VALID, constants.TIPS_COMMON_PARAM_NOT_VALID, nil))
}

// 预览
func PreviewFile(c *gin.Context) {
	var FileIdReq FileIdReq
	if c.ShouldBindQuery(&FileIdReq) == nil {
		uid, _ := c.Get("uid")
		fileInfo := models.GetFile(FileIdReq.Fid, uid.(string))
		if fileInfo == nil {
			c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_CREATE_DIR_PARAM_NOT_VALID, constants.TIPS_FILE_NOT_EXIST, nil))
			return
		}
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", "inline")
		c.Header("Content-Transfer-Encoding", "binary")
		c.File(fileInfo.FileRealPath)
		return
	}
	c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_PARAMS_NOT_VALID, constants.TIPS_COMMON_PARAM_NOT_VALID, nil))
}

// 重命名
func RenameFile(c *gin.Context) {

}

func DeleteFile(c *gin.Context) {
	var fileIdReq FileIdReq
	if c.ShouldBind(&fileIdReq) == nil {
		uid, _ := c.Get("uid")
		file := models.GetFile(fileIdReq.Fid, uid.(string))
		if file != nil {
			succ := models.DeleteFile(fileIdReq.Fid, uid.(string))
			if succ {
				// 删除实际文件
				err := os.Remove(file.FileRealPath)
				if err != nil {
					c.JSON(http.StatusInternalServerError, utils.ReturnJSON(
						constants.CODE_UNHANDLED_ERROR,
						err.Error(),
						nil,
					))
					return
				}
				c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_OK, "", nil))
				return
			}
		} else {
			c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_FILE_NOT_EXIST, constants.TIPS_FILE_NOT_EXIST, nil))
			return
		}

	} else {
		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_PARAMS_NOT_VALID, constants.TIPS_COMMON_PARAM_NOT_VALID, nil))
	}
}

// 分享资源
func ShareFileToGroup(c *gin.Context) {

}
