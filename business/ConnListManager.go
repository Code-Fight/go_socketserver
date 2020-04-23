package business

import (
	"go_socketserver/Common"
	"go_socketserver/socket"
	"net"
	"strings"
	"sync"
)

// 关闭无用的连接
// 从session中删除前先关闭连接，然后调用改方法，
// 通过遍历所有session 删除掉已经连接关闭的
// 同时产出ip对应zbm的map
func ClientListDel(conn net.Conn) {

	if conn==nil{
		//如果传入的conn不存在 那么执行删除空连接的方式
		Common.ClientList.Range(func(zbm, zbmVal interface{}) bool {
			room,zbmOk := zbmVal.(*sync.Map)

			if !zbmOk {
				return false
			}
			room.Range(func(key, value interface{}) bool {
				c,ok:=value.(*socket.Conn)

				if ok {

					if c.RECVConn ==nil || c.CMDConn==nil{
						defer room.Delete(key)


						//处理异常链接 通知下线
						if c.CMDConn !=nil{
							c.CMDConn.Close()
						}
						if c.RECVConn!=nil{
							c.RECVConn.Close()
						}


						return false

					}



				}


				return true
			})
			return true




		})
	}else {
		clientInfo:=conn.RemoteAddr().String()

		ZBM,zbmOK:=Common.ConnListIp.Load(clientInfo)

		if !zbmOK{
			return
		}

		zbmString,_:=ZBM.(string)
		room,roomOk:=Common.ClientList.Load(zbmString)
		if !roomOk {
			return
		}
		roomObj,_:=room.(*sync.Map)
		roomObj.Range(func(key, value interface{}) bool {
			c,ok := value.(*socket.Conn)

			if ok {

				if c.RECVConn ==nil || c.CMDConn==nil{
					defer roomObj.Delete(key)


					//处理异常链接 通知下线
					if c.CMDConn !=nil{
						c.CMDConn.Close()
					}
					if c.RECVConn!=nil{
						c.RECVConn.Close()
					}

					if c.CMDConn !=nil{
						// 删除ConnListIp的对应关系
						Common.ConnListIp.Delete(c.CMDConn.RemoteAddr())
						// 删除IP和通道对应关系
						Common.ConnType.Delete(c.CMDConn.RemoteAddr())
					}

					if c.RECVConn!=nil{
						// 删除ConnListIp的对应关系
						Common.ConnListIp.Delete(c.RECVConn.RemoteAddr())
						// 删除IP和通道对应关系
						Common.ConnType.Delete(c.RECVConn.RemoteAddr())

					}





					//下线通知
					statusBytes := GenDevStatusBytes(uint32(c.DevId),uint32(c.DevType),0,strings.Split(conn.RemoteAddr().String(),":")[0])



					socket.SendToRoom(0x0000, Common.Cmd_Net_Comm_Status,statusBytes, Common.SENDSOCKET,0,zbmString)
					return true
				}
				if c.RECVConn.RemoteAddr().String()==clientInfo||c.CMDConn.RemoteAddr().String()==clientInfo{
					defer roomObj.Delete(key)
					//处理异常链接 通知下线
					if c.CMDConn !=nil{
						c.CMDConn.Close()

					}
					if c.RECVConn!=nil{
						c.RECVConn.Close()
						c.RECVConn = nil
					}
					if c.CMDConn !=nil{
						// 删除ConnListIp的对应关系
						Common.ConnListIp.Delete(c.CMDConn.RemoteAddr())
						// 删除IP和通道对应关系
						Common.ConnType.Delete(c.CMDConn.RemoteAddr())
					}

					if c.RECVConn!=nil{
						// 删除ConnListIp的对应关系
						Common.ConnListIp.Delete(c.RECVConn.RemoteAddr())
						// 删除IP和通道对应关系
						Common.ConnType.Delete(c.RECVConn.RemoteAddr())

					}

					//下线通知
					statusBytes := GenDevStatusBytes(uint32(c.DevId),uint32(c.DevType),0,strings.Split(conn.RemoteAddr().String(),":")[0])
					socket.SendToRoom(0x0000, Common.Cmd_Net_Comm_Status,statusBytes, Common.SENDSOCKET,0,zbmString)

					return false
				}


			}
			return true
		})
	}

}






