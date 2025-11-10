package httpproxyrouter

import (

	"github.com/gin-gonic/gin"
)


func InitHttpProxyRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	router := gin.Default()
	router.Use(middlewares...)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong hahah",
		})
	})

	return router
}
