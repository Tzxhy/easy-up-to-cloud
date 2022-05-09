package controllers

import (
	"encoding/hex"
	"log"
	"mime"
	"net/http"
	"net/url"
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
	err = c.ShouldBind(&uploadFileReq)
	if err != nil {
		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_PARAMS_NOT_VALID, constants.CODE_PARAMS_NOT_VALID_TIPS.Tip, nil))
		return
	}
	if uploadFileReq.ParentDid != nil {
		did = *uploadFileReq.ParentDid
	}

	file := models.GetFileByName(myFile.Filename, uid.(string), did)
	if file != nil {
		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_FILENAME_HAS_BEEN_USED, constants.CODE_FILENAME_HAS_BEEN_USED_TIPS.Tip, nil))
		return
	}
	fileNameHex := hex.EncodeToString([]byte(myFile.Filename)) + filepath.Ext(myFile.Filename)
	filePath := filepath.Join(constants.UPLOAD_PATH, uid.(string), did, fileNameHex)
	utils.MakeSurePathExists(filepath.Dir(filePath))
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
			c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_CREATE_DIR_PARAM_NOT_VALID, constants.CODE_FILE_NOT_EXIST_TIPS.Tip, nil))
			return
		}
		contentType := mime.TypeByExtension(filepath.Ext(fileInfo.FileRealPath))
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		c.Header("Content-Type", contentType)
		c.Header("Content-Disposition", "attachment; filename=\""+url.QueryEscape(fileInfo.Filename)+"\"")
		c.Header("Content-Transfer-Encoding", "binary")
		c.File(fileInfo.FileRealPath)
		return
	}
	c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_PARAMS_NOT_VALID, constants.CODE_PARAMS_NOT_VALID_TIPS.Tip, nil))
}

// 预览
func PreviewFile(c *gin.Context) {
	var FileIdReq FileIdReq
	if c.ShouldBindQuery(&FileIdReq) == nil {
		uid, _ := c.Get("uid")
		fileInfo := models.GetFile(FileIdReq.Fid, uid.(string))
		if fileInfo == nil {
			c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_CREATE_DIR_PARAM_NOT_VALID, constants.CODE_FILE_NOT_EXIST_TIPS.Tip, nil))
			return
		}
		contentType := mime.TypeByExtension(filepath.Ext(fileInfo.FileRealPath))
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		c.Header("Content-Type", contentType)
		c.Header("Content-Disposition", "inline;filename=\""+url.QueryEscape(fileInfo.Filename)+"\"")
		c.Header("Content-Transfer-Encoding", "binary")
		c.File(fileInfo.FileRealPath)
		return
	}
	c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_PARAMS_NOT_VALID, constants.CODE_PARAMS_NOT_VALID_TIPS.Tip, nil))
}

type FileRenameReq struct {
	Fid  string `json:"fid" form:"fid" binding:"required"`
	Name string `json:"name" form:"name" binding:"required"`
}

// 重命名
func RenameFile(c *gin.Context) {
	var fileRenameReq FileRenameReq
	if c.ShouldBind(&fileRenameReq) != nil {
		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_PARAMS_NOT_VALID, constants.CODE_PARAMS_NOT_VALID_TIPS.Tip, nil))
		return
	}
	uid, _ := c.Get("uid")
	succ := models.RenameFile(uid.(string), fileRenameReq.Fid, fileRenameReq.Name)
	if succ {
		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_OK, "", nil))
	} else {
		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_RENAME_FILE_WITH_ERROR, constants.CODE_RENAME_FILE_WITH_ERROR_TIPS.Tip, nil))
	}
}

type MoveFileReq struct {
	Fid          string `json:"fid" form:"fid" binding:"required"`
	NewParentDid string `json:"new_parent_did" form:"new_parent_did"`
}

func MoveFile(c *gin.Context) {
	var moveFileReq MoveFileReq
	if c.ShouldBind(&moveFileReq) != nil {
		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_PARAMS_NOT_VALID, constants.CODE_PARAMS_NOT_VALID_TIPS.Tip, nil))
		return
	}
	uid, _ := c.Get("uid")
	succ := models.MoveFile(uid.(string), moveFileReq.Fid, moveFileReq.NewParentDid)
	if succ {
		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_OK, "", nil))
	} else {
		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_MOVE_FILE_WITH_ERROR, constants.CODE_MOVE_FILE_WITH_ERROR_TIPS.Tip, nil))
	}
}

func DeleteFile(c *gin.Context) {
	var fileIdReq FileIdReq
	if c.ShouldBind(&fileIdReq) == nil {
		uid, _ := c.Get("uid")
		deleteCode := deleteFile(fileIdReq.Fid, uid.(string))
		if deleteCode == OK {
			c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_OK, "", nil))
			return
		} else if FILE_NOT_FOUND == deleteCode {
			c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_FILE_NOT_EXIST, constants.CODE_FILE_NOT_EXIST_TIPS.Tip, nil))
			return
		}
		// file := models.GetFile(fileIdReq.Fid, uid.(string))
		// if file != nil {
		// 	succ := models.DeleteFile(fileIdReq.Fid, uid.(string))
		// 	if succ {
		// 		// 删除实际文件
		// 		err := os.Remove(file.FileRealPath)
		// 		if err != nil {
		// 			c.JSON(http.StatusInternalServerError, utils.ReturnJSON(
		// 				constants.CODE_UNHANDLED_ERROR,
		// 				err.Error(),
		// 				nil,
		// 			))
		// 			return
		// 		}
		// 		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_OK, "", nil))
		// 		return
		// 	}
		// } else {
		// 	c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_FILE_NOT_EXIST, constants.TIPS_FILE_NOT_EXIST, nil))
		// 	return
		// }

	} else {
		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_PARAMS_NOT_VALID, constants.CODE_PARAMS_NOT_VALID_TIPS.Tip, nil))
	}
}

const (
	INTERNAL_ERROR = iota
	FILE_NOT_FOUND
	DELETE_ERROR
	OK
)

func deleteFile(fid, uid string) uint8 {
	file := models.GetFile(fid, uid)
	if file != nil {
		succ := models.DeleteFile(fid, uid)
		if succ {
			// 删除实际文件
			err := os.Remove(file.FileRealPath)
			if err != nil {
				return INTERNAL_ERROR
			}
			return OK
		} else {
			return DELETE_ERROR
		}
	} else {
		return FILE_NOT_FOUND
	}
}

// 分享资源
func ShareFileToGroup(c *gin.Context) {

}
