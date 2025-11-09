package router

import (
	"context"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

var (
	HttpSrvHandler *http.Server
)

func HttpServerRun() {
	gin.SetMode(lib.ConfBase.DebugMode)
	r := InitRouter()
	HttpSrvHandler = &http.Server{
		Addr:           lib.GetStringConf("base.http.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(lib.GetIntConf("base.http.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(lib.GetIntConf("base.http.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(lib.GetIntConf("base.http.max_header_bytes")),
	}
	go func() {
		// 为什么不用 r.Run()？r.Run() 是对 http.ListenAndServe 的封装，但无法自定义超时等参数；用 http.Server 结构体可以更精细地配置服务。
		log.Printf(" [INFO] HttpServerRun:%s\n",lib.GetStringConf("base.http.addr"))
		if err := HttpSrvHandler.ListenAndServe(); err != nil {
			log.Fatalf(" [ERROR] HttpServerRun:%s err:%v\n", lib.GetStringConf("base.http.addr"), err)
		}
	}()
}

func HttpServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpSrvHandler.Shutdown(ctx); err != nil {
		log.Fatalf(" [ERROR] HttpServerStop err:%v\n", err)
	}
	log.Printf(" [INFO] HttpServerStop stopped\n")
}
