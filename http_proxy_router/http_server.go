package httpproxyrouter

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/puoxiu/gogate/middleware"
)

var (
	HttpSrvHandler *http.Server
)

func HttpProxyServerRun() {
	gin.SetMode(lib.ConfBase.DebugMode)
	r := InitHttpProxyRouter(
		// 全局中间件
		middleware.RequestLog(),
		middleware.RecoveryMiddleware(),
	)
	HttpSrvHandler = &http.Server{
		Addr:           lib.GetStringConf("proxy.http.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(lib.GetIntConf("proxy.http.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(lib.GetIntConf("proxy.http.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(lib.GetIntConf("proxy.http.max_header_bytes")),
	}
	log.Printf("地址是:%s", lib.GetStringConf("proxy.http.addr"))

	// 为什么不用 r.Run()？r.Run() 是对 http.ListenAndServe 的封装，但无法自定义超时等参数；用 http.Server 结构体可以更精细地配置服务。
	log.Printf(" [INFO] HttpProxyServerRun:%s\n", lib.GetStringConf("proxy.http.addr"))
	if err := HttpSrvHandler.ListenAndServe(); err != nil {
		log.Fatalf(" [ERROR] HttpProxyServerRun:%s err:%v\n", lib.GetStringConf("proxy.http.addr"), err)
	}
}

func HttpProxyServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpSrvHandler.Shutdown(ctx); err != nil {
		log.Fatalf(" [ERROR] HttpProxyServerStop err:%v\n", err)
	}
	log.Printf(" [INFO] HttpProxyServerStop stopped\n")
}
