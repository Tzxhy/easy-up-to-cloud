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
		fs.GET("get-dir-list", controllers.GetDirList)
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

	return r
}
