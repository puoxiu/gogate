package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()
	// 根路径返回当前端口，方便区分
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello from 2003 port!")
	})
	// 任意路径都返回（适配网关可能的路径剥离）
	router.NoRoute(func(c *gin.Context) {
		c.String(http.StatusOK, "Hello from 2003 port! Path: %s", c.Request.URL.Path)
	})
	// 监听 2003 端口
	router.Run(":2003")
}