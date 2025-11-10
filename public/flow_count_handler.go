package public

import (
	"sync"
	"time"
)

// FlowCounterHandler 是全局流量计数器处理器实例
var FlowCounterHandler *FlowCounter

type FlowCounter struct {
	RedisFlowCountMap map[string]*RedisFlowCountService
	RedisFlowCountSlice []*RedisFlowCountService
	Locker sync.RWMutex
}

func NewFlowCounter() *FlowCounter {
	return &FlowCounter{
		RedisFlowCountMap: map[string]*RedisFlowCountService{},
		RedisFlowCountSlice: []*RedisFlowCountService{},
		Locker: sync.RWMutex{},
	}
}

func init() {
	FlowCounterHandler = NewFlowCounter()
}

// GetCounter 获取指定服务名称的流量计数器，如果不存在则创建新的
func (counter *FlowCounter) GetCounter(serverName string) (*RedisFlowCountService, error) {
	// 遍历切片查找已存在的计数器 如果找到匹配的服务名，直接返回
	for _, item := range counter.RedisFlowCountSlice {
		if item.AppID == serverName {
			return item, nil
		}
	}

	// 如果未找到，创建新的流量计数服务实例，设置时间间隔为1秒
	newCounter := NewRedisFlowCountService(serverName, 1*time.Second)
	counter.RedisFlowCountSlice = append(counter.RedisFlowCountSlice, newCounter)
	counter.Locker.Lock()
	defer counter.Locker.Unlock()
	counter.RedisFlowCountMap[serverName] = newCounter
	return newCounter, nil
}
