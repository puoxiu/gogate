
## HTTP代理

1. 添加服务

参数说明：
1. 基本信息
* service_name：服务名称
* service_desc：服务描述

2. 路由匹配规则
* rule_type：规则类型，0：前缀URL匹配，1：域名匹配
* rule：规则，前缀URL匹配时，为URL前缀，域名匹配时，为域名
* need_https：是否需要HTTPS，0：否，1：是
* need_strip_uri：是否剥离匹配的前缀，0：否，1：是；    ----- 例如：前缀匹配为"/test_http_sss"，请求URL为"/test_http_sss/abc"，则剥离后的URL为"/abc"
* need_websocket：是否需要WebSocket，0：否，1：是

3. 路径与头信息处理
* url_rewrite：URL重写规则，正则表达式
* header_transfor：头信息转换规则，每个规则参数用空格分隔，类型只能是add、edit、del  ---- 例如：add X-Gateway gateway,edit X-Real-IP 127.0.0.1,del X-Secret

4. 权限与限流
* open_auth:  是否开启认证 开启后会校验黑白名单，0：否，1：是
* black_list：黑名单，逗号分隔的IP列表               ---- 当开启认证且白名单为空时，黑名单生效
* white_list：白名单，逗号分隔的IP列表               ---- 当开启认证且白名单不为空时，白名单生效，优先级高于黑名单
* clientip_flow_limit：客户端IP流量限制，0：不限制  
* service_flow_limit：服务端流量限制，0：不限制

5. 负载均衡
* round_type：负载均衡算法， 0=random 1=round-robin 2=weight_round-robin 3=ip_hash
* ip_list：后端服务节点列表（逗号分隔，格式 IP:端口）请求会转发到这些节点
* weight_list：节点权重列表（逗号分隔），仅在加权轮询算法下有效

6. 超时设置
* upstream_connect_timeout：与后端建立连接的超时时间，单位秒，0：默认
* upstream_header_timeout：获取后端响应头的超时时间
* upstream_idle_timeout：连接空闲的最大时间
* upstream_max_idle：最大空闲连接数

2. eg：
添加一个HTTP代理服务，以前缀匹配方式转发；负载均衡以加权轮询方法，转发到两个后端节点"127.0.0.1:2003"和"127.0.0.1:2004"，权重分别为50和50--说明平等机会
后端需运行对应的服务：
```yaml
{
  "service_name": "test_http4",
  "service_desc": "测试http代理--前缀匹配",
  "rule_type": 0,
  "rule": "/test_http_ss",
  "need_https": 0,
  "need_strip_uri": 0,
  "need_websocket": 0,
  "url_rewrite": "^/test_http_ss/(.*) $1",
  "header_transfor": "",
  "open_auth": 1,
  "black_list": "192.168.1.100,10.0.0.5",
  "white_list": "127.0.0.1,192.168.1.1, ::1",
  "clientip_flow_limit": 100,
  "service_flow_limit": 1000,
  "round_type": 1,
  "ip_list": "127.0.0.1:2003,127.0.0.1:2004",
  "weight_list": "50,50",
  "upstream_connect_timeout": 0,
  "upstream_header_timeout": 0,
  "upstream_idle_timeout": 0,
  "upstream_max_idle": 0
}
```

2. 查看服务
调用对应接口GET方法即可查看服务列表

3. 访问服务
调用对应接口POST方法即可访问服务
```bash
http://127.0.0.1:8080/test_http_ss
```
可以得到2003和2004的响应，说明负载均衡生效！



