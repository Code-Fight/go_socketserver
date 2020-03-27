package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	log "github.com/Code-Fight/golog"
	"io"
	"net"
	"os"
	"socketserver/Common"
	"socketserver/units"
	"time"

)
var ip = flag.String("ip", "127.0.0.1:2048", "ip")
var t = flag.Int("t", 10, "并发数量")
var d = flag.Int("d", 10, "并发数量")
func main() {

	flag.Parse()


	//初始化日志
	//关闭日志压缩
	//设置日志分割大小为30M
	log.Init("./log/server.log",log.DebugLevel,false,log.SetCaller(true),log.SetMaxFileSize(1024),log.SetCompress(false),log.SetMaxBackups(10))


	fmt.Println("服务器地址：",*ip)
	fmt.Println("并发客户端：",*t)

	for i := 0; i < *t; i++ {

		time.Sleep(time.Millisecond*time.Duration(*d))
		go start(*ip)
	}
	for{
		time.Sleep(time.Second*5);
	}

}

func start(ip string)  {

	defer func() {
		if err:= recover(); err != nil {
			fmt.Println(err)
		}
	}()

	idchan :=make(chan []byte,1)
	conn, err := net.Dial("tcp", ip)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		//os.Exit(1)
	}

	fmt.Println(" HandelRECVConn connect success")
	go HandelRECVConn(conn,idchan)

	conn1, err1 := net.Dial("tcp", ip)

	if err1 != nil {
		fmt.Printf("Fatal error: %s", err1.Error())
		//os.Exit(1)
	}

	fmt.Println(" HandelRECVConn connect success")
	go HandelCMDConn(conn1,idchan)

}

//发送通道
func HandelCMDConn(conn net.Conn,id <-chan []byte)  {
	//先发送登录和注册
	clientid:=<-id
	login:=[]byte{0x2a,0x00,0x00,0x00,0x0c,0x00,0x00,0xff,0xff,0x00,0x00,0x00,0x10,0x00,0x00,0x00,0x02,0x00,0x08}
	login[7] = clientid[0]
	login[8] = clientid[1]
	conn.Write(login)

	time.Sleep(time.Second*2)

	go func() {
		for{
			time.Sleep(time.Second*60)
			buff :=make([]byte,400*1024/10+100)
			binary.LittleEndian.PutUint64(buff,uint64( time.Now().UnixNano()))
			//加一堆扰乱数据
			//buff = append(buff, []byte{0x08,0x00,0x00,0x00,0x00,0x10,0x00,0x00,0x01,0x00,0x00,0x00,0x31,0x00,0x32,0x00,0x37,0x00,0x2e,0x00,0x30,0x00,0x2e,0x00,0x30,0x00,0x2e,0x00,0x31,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x04,0xa0,0x69,0x5e,0x00,0x00,0x00,0x00}...)

			for i := 0; i < 400*1024/10; i++ {
				buff = append(buff,[]byte{0x08,0x00,0x00,0x00,0x00,0x10,0x00,0x00,0x01,0x00}...)
			}



			if units.BytesToSrc(clientid) -1>=uint16(4096){
				conn.Write(Common.Packet(uint32(len(buff)+10),units.BytesToSrc(clientid),units.BytesToDes(clientid)-1,units.BytesToCmd([]byte{0x11,0x11}),uint32(len(buff)),buff))

			}
		}
	}()

	// 由于 UI客户是有CMD通道进行接收数据，所以这样进行测试。
	go func() {

		buffer := make([]byte, 4096)
		var tmpBuffer []byte
		//声明一个管道用于接收解包的数据
		readerChannel := make(chan []byte, 16)
		idchan :=make(chan []byte, 1)
		//实际业务
		go reader(conn, readerChannel, 0,idchan)

		//尝试接收数据
		for {

			n, err := conn.Read(buffer)

			if err != nil {
				if err == io.EOF {
					continue
				}
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {

					//log.Errorf("TimeOut close client: %s:",opErr.Addr.String())

				}
				//log.Error(conn.RemoteAddr().String(), " connection error: ", err, reflect.TypeOf(err))


				return
			}
			// 从缓冲区中读取数据 并尝试解包
			tmpBuffer = Common.Unpack(append(tmpBuffer, buffer[:n]...), readerChannel)
		}
	}()




}

//接收通道
func HandelRECVConn(conn net.Conn,idchan chan []byte)  {

	//先发送登录和注册
	conn.Write([]byte{0x2a,0x00,0x00,0x00,0x0c,0x00,0x00,0xff,0xff,0x00,0x00,0x00,0x10,0x00,0x00,0x00,0x02,0x00,0x08})

	//声明一个管道用于接收解包的数据
	readerChannel := make(chan []byte, 16)
	//实际业务
	go reader(conn, readerChannel, 0,idchan)

	buffer := make([]byte, 4096)
	var tmpBuffer []byte

	//尝试接收数据
	for {

		n, err := conn.Read(buffer)

		if err != nil {
			if err == io.EOF {
				continue
			}
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {

				//log.Errorf("TimeOut close client: %s:",opErr.Addr.String())

			}
			//log.Error(conn.RemoteAddr().String(), " connection error: ", err, reflect.TypeOf(err))


			return
		}
		// 从缓冲区中读取数据 并尝试解包
		tmpBuffer = Common.Unpack(append(tmpBuffer, buffer[:n]...), readerChannel)
	}


}



func reader(conn net.Conn, readerChannel <-chan []byte, timeout int,idchan chan<- []byte) {
	for {
		select {
		case data := <-readerChannel:
			p:=Common.MyProtocol{}
			p.Decode(data)
			//fmt.Printf("收到命令：%x \r\n",units.BytesToCmd(p.Data.Cmd))
			if units.BytesToCmd(p.Data.Cmd) == Common.Cmd_Login_Reply {
				//登录成功 发送id
				fmt.Println("登录成功")
				idchan<-p.Data.Des
			}

			//0x1111 随便定义的一自定义测试命令
			if  units.BytesToCmd(p.Data.Cmd) ==units.BytesToCmd([]byte{0x11,0x11}){
				//解析data 然后对比时间
				temp:= p.Data.Data[0:8]
				c:=time.Now().UnixNano()-int64(binary.LittleEndian.Uint64(temp))
				log.Infof("%v \r\n",float64(c)/1000000)
				//fmt.Printf()

			}

			break
			//case <-time.After(time.Duration(timeout) * time.Second):
			//	Log("一直没有收到数据:connection is closed."+conn.RemoteAddr().String())
			//	conn.Close()

		}
	}
}