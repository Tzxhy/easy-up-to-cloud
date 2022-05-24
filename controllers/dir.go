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
	if c.ShouldBind(&newDirInfo) != nil {
		utils.ReturnParamNotValid(c)
		return
	}
	uid, _ := c.Get("uid")
	parentDirId := constants.DIR_ROOT_ID
	if newDirInfo.ParentDirId != nil && *newDirInfo.ParentDirId != "" {
		parentDirId = *newDirInfo.ParentDirId
	}
	did, err := models.AddDir(uid.(string), newDirInfo.Name, parentDirId)
	if err == nil {
		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_OK, "", &gin.H{
			"did": did,
		}))
	} else {
		c.JSON(http.StatusOK, utils.ReturnJSON(
			constants.CODE_CREATE_DIR_WITH_ERROR_TIPS.Code,
			constants.CODE_CREATE_DIR_WITH_ERROR_TIPS.Tip,
			nil,
		))
	}
}

type GetDirInfo struct {
	DirId *string `form:"dir_id"`
}

// 获取目录信息
func GetDir(c *gin.Context) {
	var getDirInfo GetDirInfo
	if c.ShouldBindQuery(&getDirInfo) != nil {
		utils.ReturnParamNotValid(c)
		return
	}

	uid, _ := c.Get("uid")
	parentDirId := constants.DIR_ROOT_ID
	if getDirInfo.DirId != nil && *getDirInfo.DirId != "" {
		parentDirId = *getDirInfo.DirId
	}
	dirInfo := models.GetDir(parentDirId, uid.(string))
	if dirInfo != nil {
		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_OK, "", &gin.H{
			"did":         dirInfo.Did,
			"dirname":     dirInfo.Dirname,
			"parent_did":  dirInfo.ParentDid,
			"create_date": dirInfo.CreateDate,
		}))
	} else {
		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_QUERY_DIR_INFO_WITH_EMPTY_RES_TIPS.Code, constants.CODE_QUERY_DIR_INFO_WITH_EMPTY_RES_TIPS.Tip, nil))
	}

}

// 获取该目录下所有子目录，文件列表
func GetDirList(c *gin.Context) {
	var getDirInfo GetDirInfo
	if c.ShouldBindQuery(&getDirInfo) != nil {
		utils.ReturnParamNotValid(c)
		return
	}
	uid, _ := c.Get("uid")
	parentDirId := constants.DIR_ROOT_ID
	if getDirInfo.DirId != nil && *getDirInfo.DirId != "" {
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
	c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_OK, "", &gin.H{
		"dirs":  returnDirs,
		"files": returnFiles,
	}))

}

type SearchFileOrDirReq struct {
	Name string `form:"name" json:"name" binding:"required"`
}

// 查找目录或者文件
func SearchFileOrDir(c *gin.Context) {
	var searchInfo SearchFileOrDirReq
	if c.ShouldBind(&searchInfo) != nil {
		utils.ReturnParamNotValid(c)
		return
	}
	name := searchInfo.Name
	uid, _ := c.Get("uid")
	dirsInfo := models.SearchDirList(uid.(string), name)
	filesInfo := models.SearchFileList(uid.(string), name)
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
	c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_OK, "", &gin.H{
		"dirs":  returnDirs,
		"files": returnFiles,
	}))

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
		// 当删除时，同时删除分享
		models.DeleteShareByDid(deleteDirInfo.DirId)
		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_OK, "", nil))
		return
	} else {
		utils.ReturnParamNotValid(c)
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

type DeleteDirAndFileListReq struct {
	DirsId  []string `form:"dids" json:"dids" binding:"required"`
	FilesId []string `form:"fileIds" json:"fileIds" binding:"required"`
}

func DeleteDirAndFileList(c *gin.Context) {
	var deleteDirAndFileListReq DeleteDirAndFileListReq
	err := c.ShouldBind(&deleteDirAndFileListReq)
	if err != nil {
		utils.ReturnParamNotValid(c)
		return
	}

	uid, _ := c.Get("uid")
	for _, dir := range deleteDirAndFileListReq.DirsId {
		deleteDir(uid.(string), dir)
	}
	for _, file := range deleteDirAndFileListReq.FilesId {
		deleteFile(file, uid.(string))
	}
	c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_OK, "", nil))
}

type RenameDirReq struct {
	Did  string `json:"did" form:"did" binding:"required"`
	Name string `json:"name" form:"name" binding:"required"`
}

// 重命名目录
func RenameDir(c *gin.Context) {
	var renameDirReq RenameDirReq
	if c.ShouldBind(&renameDirReq) != nil {
		utils.ReturnParamNotValid(c)
		return
	}
	uid, _ := c.Get("uid")
	ok := models.RenameDir(uid.(string), renameDirReq.Did, renameDirReq.Name)
	if ok {
		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_OK, "", nil))
	} else {
		c.JSON(http.StatusOK, utils.ReturnJSON(
			constants.CODE_RENAME_DIR_WITH_ERROR,
			constants.CODE_RENAME_DIR_WITH_ERROR_TIPS.Tip,
			nil,
		))
	}
}

type MoveDirReq struct {
	Did          string `json:"did" form:"did" binding:"required"`
	NewParentDid string `json:"new_parent_did" form:"new_parent_did"`
}

// 移动文件夹
func MoveDir(c *gin.Context) {
	var moveDirReq MoveDirReq
	if c.ShouldBind(&moveDirReq) != nil {
		utils.ReturnParamNotValid(c)
		return
	}

	uid, _ := c.Get("uid")
	ok := models.MoveDir(uid.(string), moveDirReq.Did, moveDirReq.NewParentDid)
	if ok {
		c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_OK, "", nil))
		return
	} else {
		c.JSON(http.StatusOK, utils.ReturnJSON(
			constants.CODE_MOVE_DIR_WITH_ERROR,
			constants.CODE_MOVE_DIR_WITH_ERROR_TIPS.Tip,
			nil,
		))
	}
}

type MoveDirsAndFilesReq struct {
	Dirs         []string `json:"dirs" form:"dirs" binding:"required"`
	FileIds      []string `json:"fileIds" form:"fileIds" binding:"required"`
	NewParentDid string   `json:"new_parent_did" form:"new_parent_did"`
}

// 移动文件夹
func MoveDirsAndFiles(c *gin.Context) {
	var moveDirsAndFilesReq MoveDirsAndFilesReq
	if c.ShouldBind(&moveDirsAndFilesReq) != nil {
		utils.ReturnParamNotValid(c)
		return
	}
	uid, _ := c.Get("uid")
	for _, dir := range moveDirsAndFilesReq.Dirs {
		models.MoveDir(uid.(string), dir, moveDirsAndFilesReq.NewParentDid)
	}
	for _, fileId := range moveDirsAndFilesReq.FileIds {
		models.MoveFile(uid.(string), fileId, moveDirsAndFilesReq.NewParentDid)
	}
	c.JSON(http.StatusOK, utils.ReturnJSON(constants.CODE_OK, "", nil))
}
