package routers

import (
	"gitee.com/tzxhy/web/controllers"
	"gitee.com/tzxhy/web/middlewares"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.Cors())
	r.Use(middlewares.FrontendFileHandler())
	v1 := r.Group("/api/v1")

	v1.GET("ping", controllers.Ping)

	// 登录
	auth := v1.Group("auth")

	{
		auth.POST("register", controllers.Register)
		auth.POST("login", controllers.Login)
		auth.POST("logout", middlewares.NeedAuth(), controllers.Logout)
	}
	// 目录，文件
	fs := v1.Group("fs")
	fs.Use(middlewares.NeedAuth())
	{
		fs.POST("create-dir", controllers.CreateDir)
		fs.POST("delete-dir", controllers.DeleteDir)
		fs.POST("move-dir", controllers.MoveDir)
		fs.POST("rename-dir", controllers.RenameDir)
		fs.GET("get-dir-list", controllers.GetDirList) // TODO 群组增加相关字段
		fs.GET("get-dir-info", controllers.GetDir)
		fs.GET("search", controllers.SearchFileOrDir)

		fs.POST("delete-dirs-files", controllers.DeleteDirAndFileList)
		fs.POST("move-dirs-files", controllers.MoveDirsAndFiles)

		fs.POST("move-file", controllers.MoveFile)
		fs.POST("rename-file", controllers.RenameFile)
		fs.POST("upload", controllers.UploadFile)
		fs.GET("download", controllers.DownloadFile)
		fs.GET("preview", controllers.PreviewFile)
		fs.POST("delete", controllers.DeleteFile)
	}

	// 资源组
	group := v1.Group("group")
	group.Use(middlewares.NeedAuth())
	{
		// 所有操作仅操作数据库，不操作实际文件
		// 所有操作仅操作数据库，不操作实际文件
		// 所有操作仅操作数据库，不操作实际文件
		group.GET("groups", controllers.GetMyGroups)      // 获取当前账号可见群组
		group.GET("search", controllers.GroupEmpty)       // 搜索
		group.POST("create-dir", controllers.GroupEmpty)  // 创建目录
		group.POST("set-account", controllers.GroupEmpty) // 设置账户分组等信息
		group.POST("share", controllers.GroupEmpty)       // 创建分享到组，可以设计有效期
		group.POST("operation", controllers.GroupEmpty)   // 操作已共享资源。重命名，移动，删除等
		group.GET("list", controllers.GroupEmpty)         // 获取某个level的目录
		group.GET("download", controllers.GroupEmpty)     // 下载文件
		group.GET("preview", controllers.GroupEmpty)      // 资源预览
	}

	return r
}
