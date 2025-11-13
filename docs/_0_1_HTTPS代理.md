

1. 测试数据
```yaml
{
  "service_name": "test_https1",
  "service_desc": "测试HTTPS代理--前缀匹配",
  "rule_type": 0,
  "rule": "/test_https_ss1",
  "need_https": 1,
  "need_strip_uri": 1,
  "need_websocket": 0,
  "url_rewrite": "^/test_https_ss1/(.*) $1",
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