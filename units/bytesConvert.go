package units

import (
	"encoding/binary"
	"unsafe"
)

//整形转换成字节
func IntToBytes(n uint16)  []byte {
	ret :=make([]byte,2)
	binary.BigEndian.PutUint16(ret,n)
	return ret
}
func Uint32ToBytes(n uint32)  []byte {
	ret :=make([]byte,4)
	binary.BigEndian.PutUint32(ret,n)
	return ret
}

//字节转换成整形
func BytesToUint32(b []byte) uint32 {

	return binary.BigEndian.Uint32(b)

}

func BytesToUint16(b []byte) uint16 {

	return binary.BigEndian.Uint16(b)

}

func StringToBytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s)) // 获取s的起始地址开始后的两个 uintptr 指针
	h := [3]uintptr{x[0], x[1], x[1]}  // 构造三个指针数组
	return *(*[]byte)(unsafe.Pointer(&h))
}

// 实现multibytetowidechar
// 先转byte 手动多加一个0x00
func StringToBytesWithWidechar(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s)) // 获取s的起始地址开始后的两个 uintptr 指针
	h := [3]uintptr{x[0], x[1], x[1]}  // 构造三个指针数组
	data := *(*[]byte)(unsafe.Pointer(&h))
	ret :=make([]byte,0)
	for _,s:=range data{
		ret = append(ret,s,0x00)
	}
	return ret

}

//字节转string 比直接强转速度快
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func BytesToCmd(b []byte) uint16 {
	if len(b)!=2{
		return 0
	}
	// 不能直接移位只有组合 会溢出
	// 参考binary 的源码 这么写 直接用加法不行
	v:= binary.BigEndian
	return v.Uint16(b)

}

func BytesToDes(b []byte) uint16 {
	if len(b)!=2{
		return 0
	}
	// 不能直接移位只有组合 会溢出
	// 参考binary 的源码 这么写 直接用加法不行
	v:= binary.BigEndian
	return v.Uint16(b)

}


func BytesToSrc(b []byte) uint16 {
	if len(b)!=2{
		return 0
	}
	// 不能直接移位只有组合 会溢出
	// 参考binary 的源码 这么写 直接用加法不行
	v:= binary.BigEndian
	return v.Uint16(b)

}
