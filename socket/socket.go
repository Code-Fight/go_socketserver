package socket

import (
	"github.com/Code-Fight/golog"
	"io"
	"net"
	"reflect"
	"socketserver/Common"
)

var BusOnEvent func(conn net.Conn,data []byte)
var ErrorOnEvent func(conn net.Conn)






func reader(conn net.Conn, readerChannel <-chan []byte, timeout int) {
	for {
		select {
		case data := <-readerChannel:
			//conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
			BusOnEvent(conn, data)
			break
		//case <-time.After(time.Duration(timeout) * time.Second):
		//	Log("一直没有收到数据:connection is closed."+conn.RemoteAddr().String())
		//	conn.Close()

			return
		}
	}
}


func HandleConnection(conn net.Conn, timeout int) {



	//声明一个临时缓冲区，用来存储被截断的数据
	var tmpBuffer []byte

	//声明一个管道用于接收解包的数据
	readerChannel := make(chan []byte, 16)

	//实际业务
	go reader(conn, readerChannel, timeout)

	//最大4M的数据
	buffer := make([]byte, 4096)

	//尝试接收数据
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				continue
			}
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {

				log.Errorf("TimeOut close client: %s:",opErr.Addr.String())

			}
			log.Error(conn.RemoteAddr().String(), " connection error: ", err, reflect.TypeOf(err))
			socketErr(conn)

			return
		}
		// 从缓冲区中读取数据 并尝试解包
		tmpBuffer = Common.Unpack(append(tmpBuffer, buffer[:n]...), readerChannel)
	}

}

//处理socket链接异常
func socketErr(conn net.Conn) {
	ErrorOnEvent(conn)
}