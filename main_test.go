package main

import (
	"net"
	"sync"
	"testing"
)

func TestTemp(t *testing.T) {
	a:= sync.Map{}
	b:= sync.Map{}

	b.Store("a",net.Conn(nil))
	a.Store(1,b)



}
