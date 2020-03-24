package socket

import (
	"errors"
	"github.com/Code-Fight/golog"
	"net"
	"socketserver/Common"
	"socketserver/units"
)


type Conn struct {
	CMDConn  net.Conn
	RECVConn net.Conn
	DevType  uint16
	DevId    uint16
}



//给 client发送数据 通过客户端的收通道
func (c Conn) WriteData(data []byte) (n int, err error) {
	if c.RECVConn==nil{
		return 0,errors.New("conn is nil")
	}
	return c.RECVConn.Write((data))
}

//给 client回复数据 通过客户端的发数据通道 只用来回复数据
func (c Conn) ClientReply(data []byte) (n int, err error) {
	if c.CMDConn==nil{
		return 0,errors.New("conn is nil")
	}
	return c.CMDConn.Write((data))
}

// 给client回复分配的id  传入net.conn
func (c Conn) ReplyDevId(conn net.Conn, devId uint16) (n int, err error) {
	d := units.IntToBytes(devId)
	//            | * |  |        len      | |checksum| |  src    | |  dst  | |  cmd    | |       len       |
	data := []byte{0x2a, 0x00, 0x00, 0x00, 0x0a, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x12, 0x00, 0x00, 0x00, 0x00}

	data[9] = d[0]
	data[10] = d[1]

	if conn==nil{
		return 0,errors.New("conn is nil")
	}
	return conn.Write(data)
}



// 发送数据给所有在线的客户端
// 通过goroutine 启动
func SendToAll(src uint16,cmd uint16,data []byte,tunnel int,devType uint16)  {

	Common.ConnList.Range(func(key, value interface{}) bool {

		client,ok :=value.(*Conn)
		clientId,_ :=key.(uint16)
		sendData := Common.Packet(uint32(len(data)+10),src,clientId,cmd,uint32(len(data)),data)

		//根据设备进行过滤发送
		if devType!=0{
			if client.DevType!=devType{
				return true
			}
		}



		if ok{
			//单独为UI客户端进行处理一下转发的接收
			//如果转发到UI客户端的 0x1fff 那么就发到客户端的recv
			//只要是群发 就要发到客户端的recv上
			//if client.DevType == Common.Dev_Type_UI  {
			//	tunnel = Common.RECVTASK
			//}

			switch tunnel {
			//CMDTASK
			case 0:
				if client.CMDConn !=nil{
					client.CMDConn.Write(sendData)
					log.Debugf("CMDConn sendAll：%x",sendData)
				}
			//RECVTASK
			case 1:
				if client.RECVConn !=nil{
					client.RECVConn.Write(sendData)
					log.Debugf("RECVConn sendAll：%x",sendData)

				}
			}

		}

		return true
	})
}

// 发送数据给客户端
func SendData(conn *net.Conn,data []byte) (n int, err error)  {
	if (*conn)==nil{
		return 0,errors.New("conn is nil")
	}
	return (*conn).Write(data)
}
