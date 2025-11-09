package router

import (
	"log"
	"github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/redis"

	"github.com/gin-gonic/gin"
	"github.com/puoxiu/gogate/controller"
	"github.com/puoxiu/gogate/middleware"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)


func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	router := gin.Default()
	router.Use(middlewares...)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong hahah",
		})
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))


	// 设置session 存储 redis
	store, err := redis.NewStore(
		10,                                  // 1. 连接池大小
		"tcp",                               // 2. 网络类型
		"127.0.0.1:6379",  // 3. Redis 地址
		"", // 4. Redis 密码
		[]byte("secret"),                    // 6. 加密密钥（字节数组，放在最后）
	)
	if err != nil {
		log.Fatalf("redis.NewStore err:%v", err)
	}
	store.Options(sessions.Options{
		MaxAge: 86400, // 24小时过期
		Path:   "/",
	})

	// 注册业务路由
	adminLoginGroup := router.Group("/admin_login")
	adminLoginGroup.Use(
		sessions.Sessions("mysession", store),	// // "mysession" 是客户端 Cookie 键名
		middleware.RecoveryMiddleware(),
		middleware.RequestLog(),
		middleware.TranslationMiddleware(),
	)
	{
		controller.AdminLoginRegisterFunc(adminLoginGroup)
	}

	// 注册 admin 路由
	adminGroup := router.Group("/admin")
	adminGroup.Use(
		sessions.Sessions("mysession", store),
		middleware.RecoveryMiddleware(),
		middleware.RequestLog(),
		middleware.SessionAuthMiddleware(),
		middleware.TranslationMiddleware(),
	)
	{
		controller.AdminRegisterFunc(adminGroup)
	}


	serviceRouter := router.Group("/service")
	serviceRouter.Use(
		sessions.Sessions("mysession", store),
		middleware.RecoveryMiddleware(),
		middleware.RequestLog(),
		middleware.SessionAuthMiddleware(),
		middleware.TranslationMiddleware())
	{
		controller.ServiceRegisterFunc(serviceRouter)
	}

	return router
}
