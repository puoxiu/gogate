package reverse_proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/puoxiu/gogate/middleware"
	"github.com/puoxiu/gogate/reverse_proxy/load_balance"
)

func NewLoadBalanceReverseProxy(c *gin.Context, lb load_balance.LoadBalance, trans *http.Transport) *httputil.ReverseProxy {
	//请求协调者--用于修改客户端请求的目标信息，使其指向负载均衡器选择的后端服务
	director := func(req *http.Request) {
		// 1. 通过负载均衡器获取下一个后端服务地址（基于负载策略，如轮询、权重等）
		nextAddr, err := lb.Get(req.URL.String())
		//todo 优化点3
		if err != nil || nextAddr=="" {
			panic("get next addr fail")
		}
		target, err := url.Parse(nextAddr)
		if err != nil {
			panic(err)
		}
		// 2. 修改请求URL，将其指向负载均衡器选择的后端服务地址
		req.URL.Scheme = target.Scheme	// 协议（http/https）
		req.URL.Host = target.Host		// 主机+端口
		// 3. 路径合并， 例如：客户端请求：http://proxy.com/api/v1/user, 实际处理的接口：http://backend.com/service/v1/user 
		// 此时，target.Path 是后端服务的基础路径（如 /service），req.URL.Path 是客户端请求的相对路径（如 /api/v1/user），需要合并为 /service/api/v1/user 才能正确访问后端接口
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path) //
		req.Host = target.Host

		// 4. 查询参数合并
		targetQuery := target.RawQuery
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "user-agent")
		}
	}

	// 可以对响应进行修改，例如添加自定义响应头、修改响应体等
	modifyFunc := func(resp *http.Response) error {
		if strings.Contains(resp.Header.Get("Connection"), "Upgrade") {
			// 如果响应头中包含 Connection: Upgrade 不做任何修改，避免破坏连接升级流程
			return nil
		}

		//todo 优化点2
		//var payload []byte
		//var readErr error
		//
		//if strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") {
		//	gr, err := gzip.NewReader(resp.Body)
		//	if err != nil {
		//		return err
		//	}
		//	payload, readErr = ioutil.ReadAll(gr)
		//	resp.Header.Del("Content-Encoding")
		//} else {
		//	payload, readErr = ioutil.ReadAll(resp.Body)
		//}
		//if readErr != nil {
		//	return readErr
		//}
		//
		//c.Set("status_code", resp.StatusCode)
		//c.Set("payload", payload)
		//resp.Body = ioutil.NopCloser(bytes.NewBuffer(payload))
		//resp.ContentLength = int64(len(payload))
		//resp.Header.Set("Content-Length", strconv.FormatInt(int64(len(payload)), 10))
		return nil
	}

	//错误回调 ：关闭real_server时测试，错误回调
	//范围：transport.RoundTrip发生的错误、以及ModifyResponse发生的错误
	errFunc := func(w http.ResponseWriter, r *http.Request, err error) {
		middleware.ResponseError(c,999,err)
	}

	// 反向代理实例：Director用于修改客户端请求，ModifyResponse用于修改响应，ErrorHandler用于处理错误
	return &httputil.ReverseProxy{Director: director, ModifyResponse: modifyFunc, ErrorHandler: errFunc}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
