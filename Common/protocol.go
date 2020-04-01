package Common

import (
	"go_socketserver/units"
)

//封包
//@param	allen	uint16	"除头部7个字节之外的所有的长度,也就是data的长度+10"
func Packet(alllen uint32,src uint16, des uint16, cmd uint16,dataLen uint32,data []byte) []byte {
	buffer := make([]byte, 0)
	// *
	buffer = append(buffer, 0x2a)
	// 数据区长度 4字节
	bufferLen := make([]byte, 4-len(units.Uint32ToBytes(alllen)))
	bufferLen = append(bufferLen, units.Uint32ToBytes(alllen)...)
	buffer = append(buffer, bufferLen...)

	// 数据校验  2字节 参照原来的代码 这里没做校验
	buffer = append(buffer, []byte{0x00, 0x00}...)

	// Src
	buffer = append(buffer, units.IntToBytes(src)...)

	// Des
	buffer = append(buffer, units.IntToBytes(des)...)

	//Cmd
	buffer = append(buffer, units.IntToBytes(cmd)...)

	//Data Len
	_dataLen := make([]byte, 4-len(units.Uint32ToBytes((dataLen))))
	_dataLen = append(_dataLen, units.Uint32ToBytes((dataLen))...)

	buffer = append(buffer,_dataLen...)

	//data
	buffer = append(buffer,data...)


	return buffer
}

//解包
func Unpack(buffer []byte, readerChannel chan []byte) []byte {
	length := len(buffer)

	var i int
	for i = 0; i < length; i = i + 1 {

		// 第一个字符是
		// * 标记
		if buffer[i] == 0x2a {
			if length-i < 7 {
				//可能被拆包了，头部太短，等待下一次处理
				break
			}
			dataLen := int(units.BytesToUint32(buffer[i+1:i+5]))
			messageLength := dataLen + 7
			if messageLength > length-i {
				//数据区长度不够  无法解析 等待下一次
				break
			}
			data := buffer[i : i+messageLength]
			readerChannel <- data
			i = messageLength+i-1

		}

	}
	// 如果首位不是* 并且 一直没找到* 那么只能丢弃掉
	if i == length {
		//log.Printf("数据丢弃:%x",buffer)
		return make([]byte, 0)

	}
	return buffer[i:]
}


