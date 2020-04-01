package business

import (
	"go_socketserver/Common"
	"go_socketserver/socket"
	"net"
	"sync"
	"testing"
)
// 测试并发分配设备ID
func TestDistributionID(t *testing.T) {
	testData :=socket.Conn{}
	testData.DevType = Common.Dev_Type_UI
	testData.ZBM = "test"

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
	Common.ClientList.Range(func(key, value interface{}) bool {
		//t.Log(key)
		vv:=value.(sync.Map)

		vv.Range(func(key1, value1 interface{}) bool {
			index ++
			return true
		})


		return true
	})
	if index!=1000 {
		t.Error("并发ID测试失败")
	}
}
