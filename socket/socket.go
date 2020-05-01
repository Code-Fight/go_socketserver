package socket

import (
	"github.com/Code-Fight/golog"
	"go_socketserver/Common"
	"go_socketserver/units"
	"io"
	"net"
	"reflect"
)

var BusOnEvent func(conn net.Conn,data []byte,closeChannel chan struct{})
var ErrorOnEvent func(conn net.Conn)






func reader(conn net.Conn, readerChannel <-chan []byte, closeChannel chan struct{}) {
	for {
		select {
		case data := <-readerChannel:
			BusOnEvent(conn, data,closeChannel)
			break
		case _,ok := <-closeChannel:
			//如果从关闭通道中收到信息，就关闭该goroutine
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
	readerChannel := make(chan []byte, 100)

	//goroutine关闭线程
	closeChannel :=make(chan struct{})

	//实际业务
	go reader(conn, readerChannel, closeChannel)


	buffer := make([]byte, 1500)

	//尝试接收数据
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				//对端关闭 对端发送了 FIN过来 请求关闭 这里注意TCP的半关闭
				log.Infof("Client Closed!")

			}
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				log.Errorf("TimeOut close client: %s:",opErr.Addr.String())
			}
			log.Error(conn.RemoteAddr().String(), " connection error: ", err, reflect.TypeOf(err))



			socketErr(conn,closeChannel)

			return
		}
		// 从缓冲区中读取数据 并尝试解包
		tmpBuffer = Common.Unpack(append(tmpBuffer, buffer[:n]...), readerChannel)
	}

}

//处理socket链接异常
func socketErr(conn net.Conn,closeChannel chan struct{}) {
	//关闭reader goroutine
	units.SafeCloseChan(closeChannel)
	ErrorOnEvent(conn)
}
