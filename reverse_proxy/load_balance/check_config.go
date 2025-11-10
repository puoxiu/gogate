package load_balance

import (
	"fmt"
	"net"
	"reflect"
	"sort"
	"time"
)

const (
	//default check setting
	DefaultCheckMethod    = 0
	DefaultCheckTimeout   = 5
	DefaultCheckMaxErrNum = 2
	DefaultCheckInterval  = 5
)
// LoadBalanceCheckConf 负载均衡检查配置---作为被观察者
type LoadBalanceCheckConf struct {
	observers    []Observer			// 各种负载均衡器--作为观察者
	confIpWeight map[string]string	// 配置的IP权重
	activeList   []string			// 活动的IP列表
	format       string				// 配置格式
}

// AddListObserver 添加一个负载均衡器--作为观察者
func (s *LoadBalanceCheckConf) AddListObserver(o Observer) {
	s.observers = append(s.observers, o)
}

// Notify 通知所有负载均衡器--更新配置
func (s *LoadBalanceCheckConf) Notify(conf []string) {
	//fmt.Println("Notify", conf)
	//更新配置时，通知监听者也更新
	s.activeList = conf
	for _, obs := range s.observers {
		obs.Update()
	}
}

// GetConf 获取当前的负载均衡检查配置
func (s *LoadBalanceCheckConf) GetConf() []string {
	confList := []string{}
	for _, ip := range s.activeList {
		weight, ok := s.confIpWeight[ip]
		if !ok {
			weight = "50" //默认weight
		}
		confList = append(confList, fmt.Sprintf(s.format, ip)+","+weight)
	}
	return confList
}

// WatchConf 监控配置，当配置发生变化时，通知所有负载均衡器--更新配置
func (s *LoadBalanceCheckConf) WatchConf() {
	//fmt.Println("watchConf")
	go func() {
		confIpErrNum := map[string]int{}
		for {
			changedList := []string{}
			// // 遍历所有节点，通过TCP连接测试是否可用
			for item, _ := range s.confIpWeight {
				conn, err := net.DialTimeout("tcp", item, time.Duration(DefaultCheckTimeout)*time.Second)
				//todo http statuscode
				if err == nil {
					conn.Close()
					if _, ok := confIpErrNum[item]; ok {
						confIpErrNum[item] = 0
					}
				}
				if err != nil {
					confIpErrNum[item]++
				}
				if confIpErrNum[item] < DefaultCheckMaxErrNum {
					changedList = append(changedList, item)
				}
			}
			// 对比：本次健康列表和上次的 activeList 是否不一样？不一样则更新
			sort.Strings(changedList)
			sort.Strings(s.activeList)
			if !reflect.DeepEqual(changedList, s.activeList) {
				s.Notify(changedList)
			}
			time.Sleep(time.Duration(DefaultCheckInterval) * time.Second)
		}
	}()
}

// NewLoadBalanceCheckConf 创建一个新的负载均衡检查配置
func NewLoadBalanceCheckConf(format string, conf map[string]string) (*LoadBalanceCheckConf, error) {
	aList := []string{}
	//默认初始化
	for item, _ := range conf {
		aList = append(aList, item)
	}
	mConf := &LoadBalanceCheckConf{format: format, activeList: aList, confIpWeight: conf}
	mConf.WatchConf()
	return mConf, nil
}