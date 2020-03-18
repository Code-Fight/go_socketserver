package Common

import (
	"sync"
)



var (
	//注册用的锁
	RegMutex sync.Mutex

	//所有的在线用户
	ConnList sync.Map
)


