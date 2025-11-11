package tcp_proxy_middleware

import (
	"github.com/puoxiu/gogate/dao"
	"github.com/puoxiu/gogate/public"
)

// TCPFlowCountMiddleware：TCP流量统计中间件
func TCPFlowCountMiddleware() func(c *TcpSliceRouterContext) {
	return func(c *TcpSliceRouterContext) {
		serverInterface := c.Get("service")
		if serverInterface == nil {
			c.Conn.Write([]byte("get service empty"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)

		//统计项 1 全站 2 服务 3 租户
		totalCounter, err := public.FlowCounterHandler.GetCounter(public.FlowTotal)
		if err != nil {
			c.Conn.Write([]byte(err.Error()))
			c.Abort()
			return
		}
		totalCounter.Increase()

		serviceCounter, err := public.FlowCounterHandler.GetCounter(public.FlowServicePrefix + serviceDetail.Info.ServiceName)
		if err != nil {
			c.Conn.Write([]byte(err.Error()))
			c.Abort()
			return
		}
		serviceCounter.Increase()
		c.Next()
	}
}
