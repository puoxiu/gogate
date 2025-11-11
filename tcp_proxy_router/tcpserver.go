package tcp_proxy_router

import (
	"context"
	"fmt"
	"log"
	// "net"

	"github.com/puoxiu/gogate/dao"
	"github.com/puoxiu/gogate/reverse_proxy"
	"github.com/puoxiu/gogate/tcp_proxy_middleware"
	"github.com/puoxiu/gogate/tcp_server"
)

// 存储所有启动的TCP服务器实例，用于后续批量关闭, 每个服务对应一个TcpServer实例
var tcpServerList = []*tcp_server.TcpServer{}

// type tcpHandler struct {
// }
// func (t *tcpHandler) ServeTCP(ctx context.Context, src net.Conn) {
// 	src.Write([]byte("tcpHandler\n"))
// }

func TcpServerRun() {
	// 1. 获取所有需要代理的 TCP 服务配置
	serviceList := dao.ServiceManagerHandler.GetTcpServiceList()
	for _, serviceItem := range serviceList {
		tempItem := serviceItem
		// 每个服务用独立goroutine启动，支持多端口并行监听
		go func(serviceDetail *dao.ServiceDetail) {
			addr := fmt.Sprintf(":%d", serviceDetail.TCPRule.Port)
			rb, err := dao.LoadBalancerHandler.GetLoadBalancer(serviceDetail)
			if err != nil {
				log.Fatalf(" [INFO] GetTcpLoadBalancer %v err:%v\n", addr, err)
				return
			}

			//构建路由及设置中间件
			router := tcp_proxy_middleware.NewTcpSliceRouter()
			router.Group("/").Use(
				tcp_proxy_middleware.TCPFlowCountMiddleware(),
				tcp_proxy_middleware.TCPFlowLimitMiddleware(),
				tcp_proxy_middleware.TCPWhiteListMiddleware(),
				tcp_proxy_middleware.TCPBlackListMiddleware(),
			)

			//构建回调handler
			routerHandler := tcp_proxy_middleware.NewTcpSliceRouterHandler(
				func(c *tcp_proxy_middleware.TcpSliceRouterContext) tcp_server.TCPHandler {
					return reverse_proxy.NewTcpLoadBalanceReverseProxy(c, rb)
				},
				router,
			)

			baseCtx := context.WithValue(context.Background(), "service", serviceDetail)
			// 每个服务用独立的TcpServer实例，支持多端口并行监听
			tcpServer := &tcp_server.TcpServer{
				Addr:    addr,
				Handler: routerHandler,
				BaseCtx: baseCtx,
			}
			tcpServerList = append(tcpServerList, tcpServer)
			log.Printf(" [INFO] tcp_proxy_run %v\n", addr)
			if err := tcpServer.ListenAndServe(); err != nil && err != tcp_server.ErrServerClosed {
				log.Fatalf(" [INFO] tcp_proxy_run %v err:%v\n", addr, err)
			}
		}(tempItem)
	}
}

func TcpServerStop() {
	for _, tcpServer := range tcpServerList {
		tcpServer.Close()
		log.Printf(" [INFO] tcp_proxy_stop %v stopped\n", tcpServer.Addr)
	}
}


