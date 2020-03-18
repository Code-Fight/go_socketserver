package business

import (
	"net"
	"socketserver/Common"
	"socketserver/socket"
)

// 关闭无用的连接
func ConnListDel(conn net.Conn) {

	Common.ConnList.Range(func(key, value interface{}) bool {
		c,ok := value.(*socket.Conn)
		if ok {
			if c.RECVConn.RemoteAddr().String()==conn.RemoteAddr().String()||
				c.CMDConn.RemoteAddr().String()==conn.RemoteAddr().String(){
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






