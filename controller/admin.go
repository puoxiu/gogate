package controller

import (
	"encoding/json"
	"fmt"

	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/puoxiu/gogate/dao"
	"github.com/puoxiu/gogate/dto"
	"github.com/puoxiu/gogate/middleware"
	"github.com/puoxiu/gogate/public"

	"github.com/gin-contrib/sessions"
)

type AdminController struct{}

func AdminRegisterFunc(group *gin.RouterGroup) {
	adminLogin := &AdminController{}
	group.GET("/admin_info", adminLogin.AdminInfo)
	group.POST("/change_pwd", adminLogin.ChangePwd)
}

func (adminlogin *AdminController) AdminInfo(c *gin.Context) {
	sess := sessions.Default(c)
	sessInfo := sess.Get(public.AdminSessionInfoKey)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), adminSessionInfo); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	//1. 读取sessionKey对应json 转换为结构体
	//2. 取出数据然后封装输出结构体

	//Avatar       string    `json:"avatar"`
	//Introduction string    `json:"introduction"`
	//Roles        []string  `json:"roles"`
	out := &dto.AdminInfoOutput{
		ID:           adminSessionInfo.ID,
		Name:         adminSessionInfo.UserName,
		LoginTime:    adminSessionInfo.LoginTime,
		Avatar:       "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
		Introduction: "I am a super administrator",
		Roles:        []string{"admin"},
	}
	middleware.ResponseSuccess(c, out)
}

func (adminlogin *AdminController) ChangePwd(c *gin.Context) {
	params := &dto.ChangePwdInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	//1. session读取用户信息到结构体 sessInfo
	//2. sessInfo.ID 读取数据库信息 adminInfo
	//3. params.password+adminInfo.salt sha256 saltPassword
	//4. saltPassword==> adminInfo.password 执行数据保存

	//session读取用户信息到结构体
	sess := sessions.Default(c)
	sessInfo := sess.Get(public.AdminSessionInfoKey)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), adminSessionInfo); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	//从数据库中读取 adminInfo
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	adminInfo := &dao.Admin{}
	adminInfo, err = adminInfo.Find(c, tx, (&dao.Admin{UserName: adminSessionInfo.UserName}))
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	//生成新密码 saltPassword
	saltPassword := public.GenSaltPassword(adminInfo.Salt, params.Password)
	adminInfo.Password = saltPassword

	//执行数据保存
	if err := adminInfo.Save(c, tx); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	middleware.ResponseSuccess(c, "")
}
