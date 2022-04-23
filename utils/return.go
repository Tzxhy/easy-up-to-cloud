package utils

import "github.com/gin-gonic/gin"

func ReturnJSON(code int, message string, data *gin.H) gin.H {
	return gin.H{
		"code":    code,
		"message": message,
		"data":    data,
	}
}
