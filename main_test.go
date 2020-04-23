package main

import (
	"go_socketserver/Common"
	"go_socketserver/socket"
	"sync"
	"testing"
	"time"
)

func TestTemp(t *testing.T) {



}

func TestPrintClients(t *testing.T) {
	// add temp data
	a:=sync.Map{}
	a.Store(1,&socket.Conn{DevId:uint16(0x0001)})
	a.Store(2,&socket.Conn{DevId:uint16(0x0002)})
	a.Store(3,&socket.Conn{DevId:uint16(0x0003)})
	Common.ClientList.Store("CSQ",&a)


	b:=sync.Map{}
	b.Store(1,&socket.Conn{DevId:uint16(0x0004)})
	b.Store(2,&socket.Conn{DevId:uint16(0x0005)})
	b.Store(3,&socket.Conn{DevId:uint16(0x0006)})
	Common.ClientList.Store("AEQ",&b)


	//print
	PrintClients()

}

func TestSyncmap(t *testing.T) {
	// add temp data
	a:=sync.Map{}
	a.Store(1,&socket.Conn{DevId:uint16(0x0001)})
	a.Store(2,&socket.Conn{DevId:uint16(0x0002)})
	a.Store(3,&socket.Conn{DevId:uint16(0x0003)})
	Common.ClientList.Store("CSQ",&a)

	go func() {
		i:=6
		for  {

			a.Store(i,&socket.Conn{DevId:uint16(0x0001)})
			i++
			if i>90000{
				return
			}
		}

	}()

	go func() {
		i :=0
		for{
			c ,_:=Common.ClientList.Load("CSQ")
			cobj,_:= c.(*sync.Map)
			cobj.Range(func(key, value interface{}) bool {
				//t.Log(key)
				return true
			})
			i++
			if i>100000{
				return
			}
		}

	}()

	time.Sleep(time.Second*10)

}