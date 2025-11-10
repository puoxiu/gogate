package load_balance

//--------- 抽象层 --------

// 抽象的被观察者
type LoadBalanceConf interface {
	AddListObserver(o Observer)	//添加一个负载均衡器的观察者
	GetConf() []string
	WatchConf()
	Notify(conf []string)		//通知所有观察者--更新配置
}

// 抽象的观察者
type Observer interface {
	Update()		//观察者得到通知后要触发的动作
}	
