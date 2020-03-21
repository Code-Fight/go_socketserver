package socket

import (
	"github.com/Code-Fight/golog"
	"io"
	"net"
	"reflect"
	"socketserver/Common"
)

var BusOnEvent func(conn net.Conn,data []byte,closeChannel <-chan struct{})
var ErrorOnEvent func(conn net.Conn)






func reader(conn net.Conn, readerChannel <-chan []byte, closeChannel <-chan struct{}) {
	for {
		select {
		case data := <-readerChannel:
			BusOnEvent(conn, data,closeChannel)
			break
		case _,ok := <-closeChannel:
			//如果从关闭通道中收到信息，就关闭该goruntine
			if !ok {
				return
			}
			return
		}
	}
}


func HandleConnection(conn net.Conn, timeout int) {



	//声明一个临时缓冲区，用来存储被截断的数据
	var tmpBuffer []byte

	//声明一个管道用于接收解包的数据
	readerChannel := make(chan []byte, 16)

	//goruntine关闭线程
	closeChannel :=make(chan struct{})

	//实际业务
	go reader(conn, readerChannel, closeChannel)

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

			//关闭reader goruntine
			close(closeChannel)

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
