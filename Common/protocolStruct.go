package Common

import (
	"bytes"
	"encoding/binary"
	"socketserver/units"
	"time"
)

// TCP通道作用 一个收 一个发
const (
	CMDTASK int = 0
	RECVTASK int = 1
)

// 设备类型
const (
	Dev_Type_YW uint16 = 0x0001
	Dev_Type_DB uint16 = 0x0010
	Dev_Type_GB uint16 = 0x0002
	Dev_Type_DX uint16 = 0x0004
	Dev_Type_UI uint16 = 0x0008
)

//命令类型
const (
	Cmd_Heard           uint16 = 0x0000
	Cmd_Reply           uint16 = 0x0001
	Cmd_Login           uint16 = 0x0010 //设备注册
	Cmd_Pwd_User        uint16 = 0x0011 //用户名密码校验
	Cmd_Login_Reply     uint16 = 0x0012 //回复设备登录
	Cmd_No_Login        uint16 = 0x0013 //找不到目的地
	Cmd_Net_Comm_Status uint16 = 0x1004 //网络设备连接状态
)

//目的地类型
const (
	Des_UI_All uint16 = 0x1fff
	Des_All    uint16 = 0xffff
	Des_GB_ALL uint16 = 0x2fff
)

type NETCOMM_STATUS struct {
	devType uint32
	devId   uint32
	status  uint32
	ip      [32]uint8
	uptm    uint64
}

func (n *NETCOMM_STATUS) Init(id uint32, dtype uint32, status uint32,ip string) {
	n.devId = id
	n.devType = dtype
	n.status = status
	n.ipTo32byte(ip)
	n.uptm = uint64(time.Now().Unix())
}
func (n *NETCOMM_STATUS)  ipTo32byte(ip string){
	b := units.StringToBytesWithWidechar(ip)
	for i:=0;i<32-len(b);i++{
		b=append(b,0x00)
	}
	n.ip = [32]byte{}

	for i,bb:=range b{
		n.ip[i]=bb
	}
}


func (n *NETCOMM_STATUS) ToBytes() []byte {
	buf := &bytes.Buffer{}
	err := binary.Write(buf, binary.LittleEndian, n)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

type MyProtocol struct {
	Heard    string
	Len      int
	Checksum int
	Data     MyData
	SourceData []byte
}

type MyData struct {
	Src  []byte
	Des  []byte
	Cmd  []byte
	Len  []byte
	Data []byte
}

func (p *MyProtocol) Decode(buffer []byte) {
	p.SourceData = buffer
	p.Heard = "*"
	p.Len = int(units.BytesToUint32(buffer[1:5]))
	p.Checksum = int(units.BytesToUint16(buffer[5:7]))
	p.Data.Src = buffer[7:9]
	p.Data.Des = buffer[9:11]
	p.Data.Cmd = buffer[11:13]
	p.Data.Len = buffer[13:17]
	p.Data.Data = buffer[17:]
}

func (p *MyProtocol) Encode() []byte {
	buffer := make([]byte, 17)
	// *
	buffer = append(buffer, []byte(p.Heard)...)
	// 数据区长度 4字节
	bufferLen := make([]byte, 4-len(units.IntToBytes(uint16(p.Len))))
	bufferLen = append(bufferLen, units.IntToBytes(uint16(p.Len))...)
	buffer = append(buffer, bufferLen...)
	// 数据校验  2字节 参照原来的代码 这里没做校验
	buffer = append(buffer, []byte{0x00, 0x00}...)

	// Src
	buffer = append(buffer, p.Data.Src...)

	// Des
	buffer = append(buffer, p.Data.Des...)

	//Cmd
	buffer = append(buffer, p.Data.Cmd...)

	//Data Len
	buffer = append(buffer, p.Data.Len...)

	buffer = append(buffer, p.Data.Data...)

	return buffer
}
