package middlewares

import (
	"github.com/gin-gonic/gin"
)

type Store = map[string]interface{}

var store = make(Store)

func GetStore() Store {
	return store
}

func Cache() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("memoCache", store)
		c.Next()
	}
}
