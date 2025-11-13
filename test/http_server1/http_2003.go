package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func web_2003() {
	router := gin.Default()
	// 根路径返回当前端口，方便区分
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello啊,路径: %s", c.Request.URL.Path)
	})

	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello from 2003 port! Path: %s", c.Request.URL.Path)
	})

	// 任意路径都返回（适配网关可能的路径剥离）
	router.NoRoute(func(c *gin.Context) {
		c.String(http.StatusOK, "Hello from 2003 port! Path: %s", c.Request.URL.Path)
	})
	// 监听 2003 端口
	router.Run(":2003")
}

func web_2004() {
	router := gin.Default()
	// 根路径返回当前端口，方便区分
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello 2004,啊,路径: %s", c.Request.URL.Path)
	})
	// 测试路径
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello 2004,啊,路径: %s", c.Request.URL.Path)
	})

	router.NoRoute(func(c *gin.Context) {
		c.String(http.StatusOK, "Hello 2004,啊,路径: %s", c.Request.URL.Path)
	})
	// 监听 2004 端口
	router.Run(":2004")
}

func main() {
	go web_2003()
	go web_2004()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
}