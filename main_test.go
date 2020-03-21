package main

import (
	"strconv"
	"testing"
)

func TestTemp(t *testing.T) {

	a,_:=strconv.ParseUint("CSQ",128,16)

	//i :=units.StringToBytesWithWidechar("192.168.102.100")

	t.Logf("%d",a)
}
