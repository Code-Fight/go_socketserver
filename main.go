package main

import (
	"fmt"
	"net"
	"os"
	"socketserver/business"
	"socketserver/socket"
	"github.com/Code-Fight/golog"
	"socketserver/units"
)

func CheckError(err error) {
	if err != nil {
		log.Errorf("Start Error: %s", err.Error())
		os.Exit(1)
	}
}

func BusOnEvent(conn net.Conn,data []byte,closeChannel <-chan struct{})  {
	log.Debug("rev data from client:"+conn.RemoteAddr().String())
	business.CMDRoute(conn,data,closeChannel)

}

// 处理socket err
// 关闭所有相关的gorottine
// 从在线列表中删除掉该设备
func ErrorOnEvent(conn net.Conn) {
	business.ConnListDel(conn)
}


func main() {
	fmt.Println("Server Start!")
	port:=units.GetPort()
	fmt.Printf("Port:%s \r\n",port)


	//初始化日志
	//关闭日志压缩
	//设置日志分割大小为30M
	log.Init("./log/server.log",log.DebugLevel,false,log.SetCaller(true),log.SetMaxFileSize(30),log.SetCompress(false))

	//初始化 server
	netListen, err := net.Listen("tcp", "0.0.0.0:"+port)
	CheckError(err)
	defer netListen.Close()
	fmt.Print("Server Running...\r\n")
	log.Info("Waiting for clients")

	socket.BusOnEvent = BusOnEvent
	socket.ErrorOnEvent = ErrorOnEvent
	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}
		log.Info(conn.RemoteAddr().String(), " tcp connect success")
		// 如果此链接超过60秒没有发送新的数据，将被关闭
		// 超时时间 这里需要注意 如果对方不发心跳 可能会被直接关闭
		// 超时已经关闭 因为目前有的socket客户端并没有遵循发送心跳
		go socket.HandleConnection(conn, 60)
	}
}
