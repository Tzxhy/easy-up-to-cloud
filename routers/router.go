package routers

import (
	"gitee.com/tzxhy/web/controllers"
	"gitee.com/tzxhy/web/middlewares"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
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
		fs.GET("get-dir-info", controllers.GetDir)
		fs.GET("search", controllers.SearchFileOrDir)
		fs.POST("upload", controllers.UploadFile)
		fs.GET("download", controllers.DownloadFile)
		fs.POST("delete", controllers.DeleteFile)
	}

	// 资源组

	return r
}
