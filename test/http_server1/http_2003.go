package main

import (
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"runtime"

	"github.com/gin-gonic/gin"
)

func web_2003() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	// 根路径返回当前端口，方便区分
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello啊,路径: %s", c.Request.URL.Path)
	})

	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello 2003,有路径处理: %s", c.Request.URL.Path)
	})

	// 任意路径都返回（适配网关可能的路径剥离）
	router.NoRoute(func(c *gin.Context) {
		c.String(http.StatusOK, "Hello 2003,默认路径处理: %s", c.Request.URL.Path)
	})
	// 监听 2003 端口
	router.Run(":2003")
}

func web_2004() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	// 根路径返回当前端口，方便区分
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello 2004,啊,路径: %s", c.Request.URL.Path)
	})
	// 测试路径
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello 2004,有路径处理: %s", c.Request.URL.Path)
	})

	router.NoRoute(func(c *gin.Context) {
		c.String(http.StatusOK, "Hello 2004,默认路径处理: %s", c.Request.URL.Path)
	})
	// 监听 2004 端口
	router.Run(":2004")
}

func main() {
	runtime.GOMAXPROCS(4)
	debug.SetMemoryLimit(4 * 1024 * 1024 * 1024) 

	go web_2003()
	go web_2004()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
}