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
	HttpsSrvHandler *http.Server
)

func HttpsProxyServerRun() {
	gin.SetMode(lib.ConfBase.DebugMode)
	r := InitHttpProxyRouter(
				// 全局中间件
		middleware.RequestLog(),
		middleware.RecoveryMiddleware(),
	)
	HttpsSrvHandler = &http.Server{
		Addr:           lib.GetStringConf("proxy.https.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(lib.GetIntConf("proxy.https.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(lib.GetIntConf("proxy.https.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(lib.GetIntConf("proxy.https.max_header_bytes")),
	}
	log.Printf("地址是:%s", lib.GetStringConf("proxy.https.addr"))

	// 为什么不用 r.Run()？r.Run() 是对 http.ListenAndServe 的封装，但无法自定义超时等参数；用 http.Server 结构体可以更精细地配置服务。
	log.Printf(" [INFO] HttpsProxyServerRun:%s\n", lib.GetStringConf("proxy.https.addr"))
	if err := HttpsSrvHandler.ListenAndServeTLS(lib.GetStringConf("proxy.https.cert_file"), lib.GetStringConf("proxy.https.key_file")); err != nil {
		log.Fatalf(" [ERROR] HttpsProxyServerRun:%s err:%v\n", lib.GetStringConf("proxy.https.addr"), err)
	}
}

func HttpsProxyServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpsSrvHandler.Shutdown(ctx); err != nil {
		log.Fatalf(" [ERROR] HttpsProxyServerStop err:%v\n", err)
	}
	log.Printf(" [INFO] HttpsProxyServerStop stopped\n")
}
