package http_proxy_middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/puoxiu/gogate/dao"
	"github.com/puoxiu/gogate/middleware"
	"github.com/puoxiu/gogate/public"
)

// 
// 网关配置的路由规则为 /test_http_str1（前缀匹配），且开启路径剥离；
// 客户端请求路径为 /test_http_str1/abc；
// 经过中间件处理后，转发到后端的路径会变为 /abc。
func HTTPStripUriMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)

		if serviceDetail.HTTPRule.RuleType==public.HTTPRuleTypePrefixURL && serviceDetail.HTTPRule.NeedStripUri==1{
			//fmt.Println("c.Request.URL.Path",c.Request.URL.Path)
			c.Request.URL.Path = strings.Replace(c.Request.URL.Path,serviceDetail.HTTPRule.Rule,"",1)
			//fmt.Println("c.Request.URL.Path",c.Request.URL.Path)
		}


		c.Next()
	}
}
