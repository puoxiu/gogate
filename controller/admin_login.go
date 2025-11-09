package controller

import (
	// "encoding/json"
	// "time"

	"encoding/json"
	"time"

	"github.com/e421083458/golang_common/lib"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/puoxiu/gogate/dao"
	"github.com/puoxiu/gogate/dto"
	"github.com/puoxiu/gogate/middleware"
	"github.com/puoxiu/gogate/public"
)

type AdminLoginController struct{}

func AdminLoginRegister(group *gin.RouterGroup) {
	adminLogin := &AdminLoginController{}
	group.POST("/login", adminLogin.AdminLogin)
	// group.GET("/logout", adminLogin.AdminLoginOut)
}


func (adminlogin *AdminLoginController) AdminLogin(c *gin.Context) {
	params := &dto.AdminLoginReq{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 1001, err)
		return
	}

	// 1. params.UserName 取得管理员信息 admininfo
	// 2. admininfo.salt + params.Password sha256 => saltPassword
	// 3. saltPassword==admininfo.password
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	admin := &dao.Admin{}
	admin, err = admin.LoginCheck(c, tx, params)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	//设置session
	sessInfo := &dto.AdminSessionInfo{
		ID:        admin.Id,
		UserName:  admin.UserName,
		LoginTime: time.Now(),
	}
	sessBts, err := json.Marshal(sessInfo)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	sess := sessions.Default(c)
	sess.Set(public.AdminSessionInfoKey, string(sessBts))
	sess.Save()

	out := &dto.AdminLoginResp{Token: admin.UserName}
	middleware.ResponseSuccess(c, out)
}

// AdminLogin godoc
// @Summary 管理员退出
// @Description 管理员退出
// @Tags 管理员接口
// @ID /admin_login/logout
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin_login/logout [get]
// func (adminlogin *AdminLoginController) AdminLoginOut(c *gin.Context) {
// 	sess := sessions.Default(c)
// 	sess.Delete(public.AdminSessionInfoKey)
// 	sess.Save()
// 	middleware.ResponseSuccess(c, "")
// }