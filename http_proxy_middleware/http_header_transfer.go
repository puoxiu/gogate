package http_proxy_middleware

import (
	// "log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/puoxiu/gogate/dao"
	"github.com/puoxiu/gogate/middleware"
)

// HTTPHeaderTransferMiddleware 功能：根据服务配置的HeaderTransfor规则，对请求头进行添加、修改、删除操作
func HTTPHeaderTransferMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)
		// 解析HeaderTransfor配置，处理请求头
		// HeaderTransfor格式：多个规则用逗号分隔，单个规则用空格分隔为3部分（操作类型 头名称 头值）
		// 例如："add X-Gateway gateway,edit X-Real-IP 127.0.0.1,del X-Secret" 这些规则在添加服务时配置 并保存到数据库中
		// 规则之间用逗号分隔，每个规则参数用空格分隔，每个规则的操作类型只能是add、edit、del
		for _, item := range strings.Split(serviceDetail.HTTPRule.HeaderTransfor, ",") {
			items := strings.Split(item, " ")
			// log.Println("修改请求头规则: ", item)
			if len(items) != 3 {
				continue
			}
			if items[0] == "add" || items[0] == "edit" {
				c.Request.Header.Set(items[1], items[2])
			}
			if items[0] == "del" {
				c.Request.Header.Del(items[1])
			}
		}
		c.Next()
	}
}
