package main

import (
	"gitee.com/tzxhy/web/initial"
	"gitee.com/tzxhy/web/routers"
)

type LoginForm struct {
	User     string `form:"user" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func main() {
	initial.InitAll()
	// user := models.GetUserById(1)
	// fmt.Print(user)
	api := routers.InitRouter()

	api.Run(":8080")
	// ex, _ := os.Executable()
	// fmt.Print(ex)
	// fmt.Print(path.Join(path.Dir(ex)))
}
