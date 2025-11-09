package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/e421083458/golang_common/lib"
	"github.com/puoxiu/gogate/router"
)

func main()  {
	lib.InitModule("./conf/dev/",[]string{"base","mysql","redis"})
	defer lib.Destroy()
	router.HttpServerRun()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	router.HttpServerStop()
}