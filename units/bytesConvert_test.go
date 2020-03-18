package units

import "testing"

func TestBytesToCmd(t *testing.T) {

	//a := uint8(2)

	tmp := []byte{0x10,0x02}
	r:= BytesToCmd(tmp)
	t.Log(r == 0x1002)

}

func TestBytesToDes(t *testing.T) {
	//a := uint8(2)
	tmp := []byte{0x10,0x00}
	r:= BytesToCmd(tmp)
	t.Log(r)

}

func TestIntToBytes(t *testing.T) {
	t.Logf("%x",IntToBytes(4096))
}

func TestBytesToUint32(t *testing.T) {
	t.Logf("%d",BytesToUint32([]byte{0x00,0x02,0x0b,0x7f}))
}