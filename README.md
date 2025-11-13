# gogate 微服务网关
> 



## 流量统计的实现
### 实现技术栈
> 通过 Redis + 原子计数 + 定时任务 实现了高性能的请求流量统计与实时 QPS（每秒请求数）计算
核心结构：RedisFlowCountService
```go
type RedisFlowCountService struct {
	AppID       string
	Interval    time.Duration		// 统计间隔：定时持久化到Redis、计算QPS的周期（如1秒、5秒）
	QPS         int64
	Unix        int64
	TickerCount int64
	TotalCount  int64
}
```
1. 原子计数：每当有请求到达时，原子增加`TickerCount`，记录当前周期内的请求数。
2. 定时任务：在每次定时周期结束时，原子读取`TickerCount`并重置，并计算当前周期内的QPS（`TickerCount / Interval`），更新到`QPS`字段。
3. 总流量统计：原子增加`TotalCount`，记录所有请求的累计数量。

### 实现要点总结
1. 使用 atomic 原子操作 保证并发安全
2. 利用 Redis Pipeline 提升批量写入性能
3. 定时器（time.Ticker）实现异步统计与QPS计算

### key设计
> 系统在 Redis 中维护两种维度的流量统计 key：按天 和 按小时
* 按天统计： flow_day_{YYYYMMDD}_{AppID}
* 按小时统计： flow_hour_{YYYYMMDDHH}_{AppID}

### 本系统的流量统计功能：
1. HTTP代理：
    * 统计全站流量
    * 统计单个服务流量

2. TCP代理



## 限流实现
> 该模块基于 Go 官方库 golang.org/x/time/rate 实现限流
每个服务（以及每个客户端 IP）对应一个独立的 rate.Limiter 实例，通过令牌桶算法控制请求速率，防止流量突发。

通过map存储，也就是说基于内存的限流, 适用于：
* 单机部署、或者每个节点独立限流
* 性能最高（没有网络开销）通信快

后期如果需要分布式限流，则可以基于redis实现





