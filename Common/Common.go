package Common

import (
	"sync"
)



var (
	//注册用的锁
	RegMutex sync.Mutex

	//连接对应的通道类型
	//为快速确定该连接的通道类型
	//在新项目中，不建议使用该项目中的所谓双通道
	//浪费链接，完全没有必要
	ConnType sync.Map

	//所有的在线用户ip和设备id的对应关系
	//为快速找到需要关闭的设备
	ConnListIp sync.Map

	// 根据车站进行区分房间的用户列表
	// key为车站码 value为某个车站所有设备map
	ClientList sync.Map
)

