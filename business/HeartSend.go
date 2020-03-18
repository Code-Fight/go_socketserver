package business

import (
	"github.com/Code-Fight/golog"
	"net"
	"socketserver/Common"
	"time"
)

func HeardEvent(conn net.Conn,dst uint16) {
	errCount := 0
	time.Sleep(time.Second * 30)

	for true {
		time.Sleep(time.Second * 30)
		//[]byte{0x00,0x00,0x00,0x00,0x00,0x10,0x00,0x00,0x00,0x00,0x00,0x00}

		sendData := Common.Packet(10,0,dst, Common.Cmd_Heard,0,nil)
		if conn==nil{
			log.Errorf("conn不存在，sendData:%x",sendData)
			return
		}

		//log.Printf("%s 心跳：%x",conn.RemoteAddr().String(),sendData)

		_,err := conn.Write(sendData)

		if err != nil{
			errCount++
		}else {
			errCount = 0
			//log.Println("发送心跳")
		}

		if errCount >3{
			conn.Close()
			//socketErr(conn)
		}
	}

}
