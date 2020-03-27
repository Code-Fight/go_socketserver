package Common

import (
	"sync"
)



var (
	//注册用的锁
	RegMutex sync.Mutex

	//所有的在线用户
	//ConnList sync.Map

	//所有的在线用户ip和设备id的对应关系
	//为快速找到需要关闭的设备
	ConnListIp sync.Map

	// 根据车站进行区分房间的用户列表
	// key为车站码 value为某个车站所有设备map
	ClientList sync.Map
)

