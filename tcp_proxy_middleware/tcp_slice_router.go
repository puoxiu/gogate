package tcp_proxy_middleware

import (
	"context"
	"math"
	"net"

	"github.com/puoxiu/gogate/tcp_server"
)

const abortIndex int8 = math.MaxInt8 / 2 //最多 63 个中间件


// TcpHandlerFunc：中间件函数类型，接收上下文并处理
type TcpHandlerFunc func(*TcpSliceRouterContext)

// TcpSliceRouter：路由核心，管理所有分组（Group）
type TcpSliceRouter struct {
	groups []*TcpSliceGroup
}

// TcpSliceGroup：分组，每个分组可以有不同的路径和中间件
type TcpSliceGroup struct {
	*TcpSliceRouter
	path     string
	handlers []TcpHandlerFunc
}

// TcpSliceRouterContext：中间件上下文，传递连接、状态、分组信息
type TcpSliceRouterContext struct {
	Conn net.Conn
	Ctx  context.Context
	*TcpSliceGroup
	index int8
}

func newTcpSliceRouterContext(conn net.Conn, r *TcpSliceRouter, ctx context.Context) *TcpSliceRouterContext {
	newTcpSliceGroup := &TcpSliceGroup{}
	*newTcpSliceGroup = *r.groups[0] //浅拷贝数组指针,只会使用第一个分组
	c := &TcpSliceRouterContext{Conn: conn, TcpSliceGroup: newTcpSliceGroup, Ctx: ctx,}
	c.Reset()
	return c
}

func (c *TcpSliceRouterContext) Get(key interface{}) interface{} {
	return c.Ctx.Value(key)
}

func (c *TcpSliceRouterContext) Set(key, val interface{}) {
	c.Ctx = context.WithValue(c.Ctx, key, val)
}

// TcpSliceRouterHandler：适配tcp_server.TCPHandler接口的处理器
type TcpSliceRouterHandler struct {
	coreFunc func(*TcpSliceRouterContext) tcp_server.TCPHandler
	router   *TcpSliceRouter
}

// NewTcpSliceRouterHandler：创建路由处理器
func NewTcpSliceRouterHandler(coreFunc func(*TcpSliceRouterContext) tcp_server.TCPHandler, router *TcpSliceRouter) *TcpSliceRouterHandler {
	return &TcpSliceRouterHandler{
		coreFunc: coreFunc,
		router:   router,
	}
}

// ServeTCP：实现 tcp_server.TCPHandler 接口
func (w *TcpSliceRouterHandler) ServeTCP(ctx context.Context, conn net.Conn) {
	c := newTcpSliceRouterContext(conn, w.router, ctx)
	c.handlers = append(c.handlers, func(c *TcpSliceRouterContext) {
		w.coreFunc(c).ServeTCP(ctx, conn)
	})
	c.Reset()
	c.Next()
}


// 构造 router
func NewTcpSliceRouter() *TcpSliceRouter {
	return &TcpSliceRouter{}
}

// 创建 Group
func (g *TcpSliceRouter) Group(path string) *TcpSliceGroup {
	if path != "/" {
		panic("only accept path=/")
	}
	return &TcpSliceGroup{
		TcpSliceRouter: g,
		path:           path,
	}
}

// Use 添加中间件, 其实就是将处理函数添加到 handlers 数组中
func (g *TcpSliceGroup) Use(middlewares ...TcpHandlerFunc) *TcpSliceGroup {
	g.handlers = append(g.handlers, middlewares...)
	existsFlag := false
	for _, oldGroup := range g.TcpSliceRouter.groups {
		if oldGroup == g {
			existsFlag = true
		}
	}
	if !existsFlag {
		g.TcpSliceRouter.groups = append(g.TcpSliceRouter.groups, g)
	}
	return g
}

// 从最先加入中间件开始回调
func (c *TcpSliceRouterContext) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

// 跳出中间件方法
func (c *TcpSliceRouterContext) Abort() {
	c.index = abortIndex
}

// 是否跳过了回调
func (c *TcpSliceRouterContext) IsAborted() bool {
	return c.index >= abortIndex
}

// 重置回调
func (c *TcpSliceRouterContext) Reset() {
	c.index = -1
}
