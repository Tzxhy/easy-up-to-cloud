package controllers

import (
	"net/http"

	"gitee.com/tzxhy/web/constants"
	"gitee.com/tzxhy/web/models"
	"gitee.com/tzxhy/web/utils"
	"github.com/gin-gonic/gin"
)

type NewDirInfo struct {
	ParentDirId *string `json:"parent_id"`
	Name        string  `json:"name" binding:"required"`
}

// 创建目录
func CreateDir(c *gin.Context) {
	var newDirInfo NewDirInfo
	err := c.ShouldBind(&newDirInfo)
	if err == nil {
		uid, _ := c.Get("uid")
		parentDirId := ""
		if newDirInfo.ParentDirId != nil {
			parentDirId = *newDirInfo.ParentDirId
		}
		did, err := models.AddDir(uid.(string), newDirInfo.Name, parentDirId)
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
	DirId *string `form:"dir_id"`
}

// 获取目录信息
func GetDir(c *gin.Context) {
	var getDirInfo GetDirInfo
	if err := c.ShouldBindQuery(&getDirInfo); err == nil {
		uid, _ := c.Get("uid")
		parentDirId := ""
		if getDirInfo.DirId != nil {
			parentDirId = *getDirInfo.DirId
		}
		dirInfo := models.GetDir(parentDirId, uid.(string))
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

// 获取该目录下所有子目录，文件列表
func GetDirList(c *gin.Context) {
	var getDirInfo GetDirInfo
	if err := c.ShouldBindQuery(&getDirInfo); err == nil {
		uid, _ := c.Get("uid")
		parentDirId := ""
		if getDirInfo.DirId != nil {
			parentDirId = *getDirInfo.DirId
		}
		dirsInfo := models.GetDirList(parentDirId, uid.(string))
		filesInfo := models.GetFileList(parentDirId, uid.(string))
		returnDirs := make([]models.Dir, 0)
		returnFiles := make([]models.File, 0)
		if dirsInfo != nil {
			if len(*dirsInfo) > 0 {
				returnDirs = *dirsInfo
			}
		}
		if filesInfo != nil {
			if len(*filesInfo) > 0 {
				returnFiles = *filesInfo
			}
		}
		// c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_QUERY_DIR_INFO_WITH_EMPTY_RES, constants.TIPS_QUERY_DIR_INFO_WITH_EMPTY_RES, nil))
		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_OK, "", &gin.H{
			"dirs":  returnDirs,
			"files": returnFiles,
		}))
	} else {
		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_PARAMS_NOT_VALID, constants.TIPS_COMMON_PARAM_NOT_VALID, nil))
	}
}

// 查找目录或者文件
func SearchFileOrDir(c *gin.Context) {

}

type DeleteDirInfo struct {
	DirId string `form:"did" json:"did" binding:"required"`
}

// 删除目录
// 会删除对应数据库，删除对应的文件
// 高危操作
func DeleteDir(c *gin.Context) {
	var deleteDirInfo DeleteDirInfo
	if c.ShouldBind(&deleteDirInfo) == nil {
		// models.Get
		uid, _ := c.Get("uid")
		deleteDir(uid.(string), deleteDirInfo.DirId)
		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_OK, "", nil))
		return
	} else {
		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_PARAMS_NOT_VALID, constants.TIPS_COMMON_PARAM_NOT_VALID, nil))
		return
	}
}

// 删除目录
func deleteDir(owner_id, did string) {
	// 先找到所有文件夹，文件，递归执行删除，再删除文件夹
	dirs := models.GetDirList(did, owner_id)
	files := models.GetFileList(did, owner_id)

	for _, item := range *files {
		deleteFile(item.Fid, owner_id)
	}

	for _, item := range *dirs {
		deleteDir(owner_id, item.Did)
	}

	models.DeleteSingleDir(owner_id, did)

}

type RenameDirReq struct {
	Did  string `json:"did" form:"did" binding:"required"`
	Name string `json:"name" form:"name" binding:"required"`
}

// 重命名目录
func RenameDir(c *gin.Context) {
	var renameDirReq RenameDirReq
	if c.ShouldBind(&renameDirReq) != nil {
		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_PARAMS_NOT_VALID, constants.TIPS_COMMON_PARAM_NOT_VALID, nil))
		return
	}
	uid, _ := c.Get("uid")
	succ := models.RenameDir(uid.(string), renameDirReq.Did, renameDirReq.Name)
	if succ {
		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_OK, "", nil))
	} else {
		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_RENAME_DIR_WITH_ERROR, constants.TIPS_RENAME_DIR_WITH_ERROR, nil))
	}
}

type MoveDirReq struct {
	Did          string `json:"did" form:"did" binding:"required"`
	NewParentDid string `json:"new_parent_did" form:"new_parent_did"`
}

// 移动文件夹
func MoveDir(c *gin.Context) {
	var moveDirReq MoveDirReq
	if c.ShouldBind(&moveDirReq) == nil {
		uid, _ := c.Get("uid")
		succ := models.MoveDir(uid.(string), moveDirReq.Did, moveDirReq.NewParentDid)
		if succ {
			c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_OK, "", nil))
			return
		} else {
			c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_MOVE_DIR_WITH_ERROR, constants.TIPS_MOVE_DIR_WITH_ERROR, nil))
		}
	} else {
		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_PARAMS_NOT_VALID, constants.TIPS_COMMON_PARAM_NOT_VALID, nil))
	}
}
