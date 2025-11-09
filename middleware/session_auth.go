package middleware


import (
	// "errors"
	"github.com/gin-gonic/gin"
)

func SessionAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// session := session.Default(c)
		// if name,ok:=session.Get("user").(string);!ok||name==""{
		// 	ResponseError(c, InternalErrorCode, errors.New("user not login"))
		// 	c.Abort()
		// 	return
		// }
		c.Next()
	}
}
