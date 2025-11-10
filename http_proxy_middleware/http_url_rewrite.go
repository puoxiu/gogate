package http_proxy_middleware

import (
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/puoxiu/gogate/dao"
	"github.com/puoxiu/gogate/middleware"
)

// URL 路径重写的中间件， 通过正则匹配对请求路径进行替换
// 配置规则为 ^/test/(.*) /$1（表示将 /test/abc 重写为 /abc）；
// 客户端请求路径为 /test/abc；
// 经过中间件处理后，转发到后端的路径会变为 /abc。
func HTTPUrlRewriteMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)
		for _,item:=range strings.Split(serviceDetail.HTTPRule.UrlRewrite,","){
			//fmt.Println("item rewrite",item)
			items:=strings.Split(item," ")
			if len(items)!=2{
				continue
			}
			regexp,err:=regexp.Compile(items[0])
			if err!=nil{
				//fmt.Println("regexp.Compile err",err)
				continue
			}
			//fmt.Println("before rewrite",c.Request.URL.Path)
			replacePath:=regexp.ReplaceAll([]byte(c.Request.URL.Path),[]byte(items[1]))
			c.Request.URL.Path = string(replacePath)
			//fmt.Println("after rewrite",c.Request.URL.Path)
		}
		c.Next()
	}
}
