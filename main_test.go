package main

import (
	"socketserver/Common"
	"socketserver/socket"
	"sync"
	"testing"
)

func TestTemp(t *testing.T) {



}

func TestPrintClients(t *testing.T) {
	// add temp data
	a:=sync.Map{}
	a.Store(1,&socket.Conn{DevId:uint16(0x0001)})
	a.Store(2,&socket.Conn{DevId:uint16(0x0002)})
	a.Store(3,&socket.Conn{DevId:uint16(0x0003)})
	Common.ClientList.Store("CSQ",a)

	b:=sync.Map{}
	b.Store(1,&socket.Conn{DevId:uint16(0x0004)})
	b.Store(2,&socket.Conn{DevId:uint16(0x0005)})
	b.Store(3,&socket.Conn{DevId:uint16(0x0006)})
	Common.ClientList.Store("AEQ",b)


	//print
	PrintClients()

}