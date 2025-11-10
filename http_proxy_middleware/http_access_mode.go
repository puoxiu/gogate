package http_proxy_middleware

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/puoxiu/gogate/dao"
	"github.com/puoxiu/gogate/middleware"
	"github.com/puoxiu/gogate/public"
)

// HTTPAccessModeMiddleware 根据请求信息（如路径、域名等）匹配对应的后端服务配置，是网关路由的基础
func HTTPAccessModeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		service, err := dao.ServiceManagerHandler.HTTPAccessMode(c)
		if err != nil {
			log.Println("没有匹配到服务===")
			middleware.ResponseError(c, 1001, err)
			c.Abort()
			return
		}
		log.Println("matched service",public.Obj2Json(service))
		c.Set("service", service)
		c.Next()
	}
}
