package initial

import (
	"gitee.com/tzxhy/web/constants"
	"gitee.com/tzxhy/web/models"
	"gitee.com/tzxhy/web/utils"
)

func InitAll() {
	utils.ShowAppInfo()
	models.InitSqlite3()
	utils.MakeSurePathExists(constants.UPLOAD_PATH)
	InitStatic()
}
