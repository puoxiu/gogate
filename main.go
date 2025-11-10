package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/e421083458/golang_common/lib"
	"github.com/puoxiu/gogate/dao"
	httpproxyrouter "github.com/puoxiu/gogate/http_proxy_router"
	"github.com/puoxiu/gogate/router"
)

var (
	conf = flag.String("conf", "", "the config path, like ./conf/dev/")
	mode       = flag.String("mode", "", "the server or dashboard")
)

func init() {
	flag.Parse()
	if *conf == "" || *mode == "" {
		flag.Usage()
		os.Exit(1)
	}
}

// go run main.go -conf ./conf/dev/ -mode server
// go run main.go -conf ./conf/dev/ -mode dashboard
func main()  {
	fmt.Println("conf:", *conf)
	lib.InitModule(*conf,[]string{"base","mysql","redis"})
	defer lib.Destroy()

	switch *mode {
	case "dashboard":
		startDashboard()
	case "server":
		startServer()
	default:
		flag.Usage()
		os.Exit(1)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	stop()
}



// startDashboard 启动后台管理服务
func startDashboard() {
	fmt.Println("======mode====", *mode)
	router.HttpServerRun()
}

// startServer 启动HTTP和HTTPS代理服务器
func startServer() {
	fmt.Println("======mode====", *mode)

	// 初始化 从 DB 加载所有 HTTP 服务配置到内存中
	dao.ServiceManagerHandler.LoadOnce()

	go func() {
		httpproxyrouter.HttpProxyServerRun()
	}()

	go func() {
		httpproxyrouter.HttpsProxyServerRun()
	}()
}

func stop() {
	switch *mode {
	case "dashboard":
		router.HttpServerStop()
	case "server":
		httpproxyrouter.HttpProxyServerStop()
		httpproxyrouter.HttpsProxyServerStop()
	default:
		flag.Usage()
		os.Exit(1)
	}
	
}