package main

import (
	"socketserver/units"
	"testing"
)

func TestTemp(t *testing.T) {
	i :=units.StringToBytesWithWidechar("192.168.102.100")

	t.Logf("%x",i)
}
