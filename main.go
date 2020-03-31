package main

import (
	"fmt"
	"github.com/Code-Fight/golog"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"socketserver/Common"
	"socketserver/business"
	"socketserver/socket"
	"socketserver/units"
	"sync"
	"time"
)

func CheckError(err error) {
	if err != nil {
		log.Errorf("Start Error: %s", err.Error())
		os.Exit(1)
	}
}

func BusOnEvent(conn net.Conn,data []byte,closeChannel chan struct{})  {
	log.Debug("rev data from client:"+conn.RemoteAddr().String())
	business.CMDRoute(conn,data,closeChannel)

}

// 处理socket err
// 关闭所有相关的gorottine
// 从在线列表中删除掉该设备
func ErrorOnEvent(conn net.Conn) {
	business.ClientListDel(conn)
}


func main() {


	fmt.Println("Server Start!")
	units.ConfigInit()
	port:=units.GetPort()
	logLevel:=units.GetLog()
	if logLevel!=log.ErrorLevel && logLevel!=log.DebugLevel  && logLevel!=log.InfoLevel  && logLevel!=log.PanicLevel && logLevel!=log.WarnLevel {
		fmt.Println("the config log level error")
		os.Exit(0)
	}

	// 如果是debug 开启一个打印方法
	if logLevel==log.DebugLevel{
		PrintClients()
	}


	//初始化日志
	//关闭日志压缩
	//设置日志分割大小为30M
	log.Init("./log/server.log",logLevel,false,log.SetCaller(true),log.SetMaxFileSize(30),log.SetCompress(false))

	go InitTcpServer(port)

	fmt.Print("HTTP Server Running on port:16060\r\n")
	httpErr:= http.ListenAndServe("0.0.0.0:16060", nil)
	if httpErr!=nil{
		fmt.Println("HTTP Server ERR:",httpErr.Error())
	}

}

func InitTcpServer(port string) {
	//初始化 server
	netListen, err := net.Listen("tcp", "0.0.0.0:"+port)
	CheckError(err)
	defer netListen.Close()
	fmt.Print("TCP Server Running on Port:",port+"\r\n")
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

func PrintClients() {
	go func() {
		for {
			fmt.Println("================[ ", time.Now().String(), " ]================")

			Common.ClientList.Range(func(roomKey, roomVal interface{}) bool {
				room,_:=roomVal.(sync.Map)
				fmt.Print("[ZBM] ",roomKey," [Clients] ")

				room.Range(func(c, cVal interface{}) bool {
					client,_:=cVal.(*socket.Conn)




					fmt.Printf("%x  |  ",units.IntToBytes(client.DevId))
					return true
				})
				fmt.Println("")
				return true
			})



			fmt.Println("")
			time.Sleep(time.Second * 5)
		}
	}()

}