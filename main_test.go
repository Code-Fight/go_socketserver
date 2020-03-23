package main

import (
	"net"
	"sync"
	"testing"
	"time"
)

func TestTemp(t *testing.T) {
	a:= sync.Map{}


	old:=time.Now().UnixNano()



	for i:=0;i<1000000;i++{


			b:= sync.Map{}
			b.Store(string(i)+"-1",net.Conn(nil))
			b.Store(string(i)+"-2",net.Conn(nil))
			a.Store(i,b)



	}

	t.Log("加载用时：",time.Now().UnixNano()-old)

	w:=sync.WaitGroup{}

	old=time.Now().UnixNano()
	for n := 0; n < 100; n++ {
		w.Add(1000000)
		for i:=0;i<1000000;i++{


			go func(i int) {
				v,ok:=a.Load(i)
				if ok{
					vv,yes:=v.(sync.Map)
					if yes{
						vv.Load(string(i)+"-1")
						vv.Load(string(i)+"-2")
						w.Done()
					}
				}
				//for j := i * 1; j < (i+1)*1; j++ {
				//	v,ok:=a.Load(i)
				//	if ok{
				//		vv,yes:=v.(sync.Map)
				//		if yes{
				//			vv.Load(string(i)+"-1")
				//			vv.Load(string(i)+"-2")
				//			w.Done()
				//		}
				//	}
				//}


			}(i)

		}
	}

	w.Wait()
	t.Log("读取用时：",(time.Now().UnixNano()-old)/100)










}
