package Common

import (
	"fmt"
	"testing"
)

func Test_NETCOMM_STATUS(t *testing.T) {
	a := NETCOMM_STATUS{}
	a.Init(1,1,1,"192.168.100.100")
	a.ToBytes()
	fmt.Printf("%x",a.ToBytes())
}
