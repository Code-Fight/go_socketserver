package socket

import (
	"socketserver/Common"
	"testing"
)

func TestTemp(t *testing.T)  {
	//a := []byte{0x00,0x00}
	//b:=


}



func TestConnList(t *testing.T)  {
	for i := 0; i < 200; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				if _,err := Common.ConnList.Load(j);!err {
					Common.ConnList.Store(j, Conn{})

				}else
				{
					Common.ConnList.Delete(j)
				}

				length := 0

				Common.ConnList.Range(func(_, _ interface{}) bool {
					length++

					return true
				})

				//t.Log(length)
			}
		}()
	}

	t.Log("TestConnList success")


}