package business

import (
	"net"
	"socketserver/Common"
	"socketserver/socket"
	"strings"
)

// 关闭无用的连接
// 从session中删除前先关闭连接，然后调用改方法，
// 通过遍历所有session 删除掉已经连接关闭的
func ConnListDel(conn net.Conn) {

	if conn==nil{
		//如果传入的conn不存在 那么执行删除空连接的方式
		Common.ConnList.Range(func(key, value interface{}) bool {
			c,ok := value.(*socket.Conn)

			if ok {

				if c.RECVConn ==nil || c.CMDConn==nil{
					defer Common.ConnList.Delete(key)


					//处理异常链接 通知下线
					if c.CMDConn !=nil{
						c.CMDConn.Close()
					}
					if c.RECVConn!=nil{
						c.RECVConn.Close()
					}

					//下线通知
					statusBytes := GenDevStatusBytes(uint32(c.DevId),uint32(c.DevType),0,strings.Split(conn.RemoteAddr().String(),":")[0])
					socket.SendToAll(0x0000, Common.Cmd_Net_Comm_Status,statusBytes, Common.RECVTASK,0)

					//关闭现有的连接
					//if c.CMDConn !=nil{
					//	cmderr:=c.CMDConn.Close()
					//	if cmderr!=nil{
					//
					//	}
					//}
					//
					//if c.RECVConn!=nil{
					//	recverr:=c.RECVConn.Close()
					//	if recverr!=nil{
					//
					//	}
					//}

					return true

				}



			}
			return true
		})
	}else {
		clientInfo:=conn.RemoteAddr().String()


		Common.ConnList.Range(func(key, value interface{}) bool {
			c,ok := value.(*socket.Conn)

			if ok {

				if c.RECVConn ==nil || c.CMDConn==nil{
					defer Common.ConnList.Delete(key)


					//处理异常链接 通知下线
					if c.CMDConn !=nil{
						c.CMDConn.Close()
					}
					if c.RECVConn!=nil{
						c.RECVConn.Close()
					}

					//下线通知
					statusBytes := GenDevStatusBytes(uint32(c.DevId),uint32(c.DevType),0,strings.Split(conn.RemoteAddr().String(),":")[0])
					socket.SendToAll(0x0000, Common.Cmd_Net_Comm_Status,statusBytes, Common.RECVTASK,0)
					return true
				}
				if c.RECVConn.RemoteAddr().String()==clientInfo||c.CMDConn.RemoteAddr().String()==clientInfo{
					defer Common.ConnList.Delete(key)
					//处理异常链接 通知下线
					if c.CMDConn !=nil{
						c.CMDConn.Close()
					}
					if c.RECVConn!=nil{
						c.RECVConn.Close()
					}
					//下线通知
					statusBytes := GenDevStatusBytes(uint32(c.DevId),uint32(c.DevType),0,strings.Split(conn.RemoteAddr().String(),":")[0])
					socket.SendToAll(0x0000, Common.Cmd_Net_Comm_Status,statusBytes, Common.RECVTASK,0)
					return false
				}


			}
			return true
		})
	}

}

func GetZBMByConn()  {

}




