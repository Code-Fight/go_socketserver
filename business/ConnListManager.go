package business

import (
	"net"
	"socketserver/Common"
	"socketserver/socket"
	"strings"
)

// 关闭无用的连接
func ConnListDel(conn net.Conn) {

	Common.ConnList.Range(func(key, value interface{}) bool {
		c,ok := value.(*socket.Conn)
		if ok {
			if c.RECVConn.RemoteAddr().String()==conn.RemoteAddr().String()||
				c.CMDConn.RemoteAddr().String()==conn.RemoteAddr().String(){
				//下线通知

				statusBytes := GenDevStatusBytes(uint32(c.DevId),uint32(c.DevType),0,strings.Split(conn.RemoteAddr().String(),":")[0])
				socket.SendToAll(0x0000, Common.Cmd_Net_Comm_Status,statusBytes, Common.RECVTASK,0)

				//关闭现有的连接
				c.CMDConn.Close()
				c.RECVConn.Close()
				defer Common.ConnList.Delete(key)
				return true
			}

		}
		return true
	})
}






