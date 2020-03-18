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

func BusOnEvent(conn net.Conn,data []byte)  {
	log.Debug("rev data from client:"+conn.RemoteAddr().String())
	business.CMDRoute(conn,data)

}

// 处理socket err
// 关闭所有相关的goruntine
// 从在线列表中删除掉该设备
func ErrorOnEvent(conn net.Conn) {
	business.ConnListDel(conn)
}


func main() {
	fmt.Println("Server Start!")
	port:=units.GetPort()
	fmt.Printf("Port:%s",port)


	//初始化日志
	log.Init("./serverlog",log.DebugLevel,false,log.SetCaller(true))

	//初始化 server
	netListen, err := net.Listen("tcp", "0.0.0.0:"+port)
	CheckError(err)
	defer netListen.Close()
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
		go socket.HandleConnection(conn, 60)
	}
}
