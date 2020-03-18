package business

import (
	"net"
	"socketserver/Common"
	"socketserver/socket"
	"sync"
	"testing"
)
// 测试并发分配设备ID
func TestDistributionID(t *testing.T) {
	testData :=socket.Conn{}
	testData.DevType = Common.Dev_Type_UI

	var wg sync.WaitGroup
	wg.Add(1000)
	for i := 0;i<1000 ;i++  {
		go func() {

			temp :=net.Conn(nil)
			distributionID(&testData,&temp)
			wg.Done()
		}()

	}
	wg.Wait()

	index := 0
	Common.ConnList.Range(func(key, value interface{}) bool {
		//t.Log(key)
		index ++
		return true
	})
	if index!=1000 {
		t.Error("并发ID测试失败")
	}
}
