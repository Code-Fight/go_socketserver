package business

import (
	"github.com/Code-Fight/golog"
	"net"
	"socketserver/Common"
	"socketserver/socket"
	"socketserver/units"
)

// 命令路由
func CMDRoute(conn net.Conn, data []byte ) {

	log.Debugf("收到%s 的数据包：%x",conn.RemoteAddr().String(),data)
	var s = Common.MyProtocol{}
	s.Decode(data)
	cmd := units.BytesToCmd(s.Data.Cmd)

	if cmd != Common.Cmd_Login {
		// 检测是否设置src
		if _,ok := Common.ConnList.Load(units.BytesToSrc(s.Data.Src));!ok{
			log.Error("设备未注册")
			conn.Close()
		}
	}


	switch cmd {
	case Common.Cmd_Heard:
		log.Debug("recv Cmd_Heard, [PASS]")

	case Common.Cmd_Reply:
		log.Debug("recv Cmd_Reply,pass [PASS]")

	case Common.Cmd_Login:
		log.Debug("recv Cmd_Login")
		Reg(conn,&s)

	case Common.Cmd_Login_Reply:
		log.Debug("recv Cmd_Login_Reply [PASS]")

	case Common.Cmd_No_Login:
		log.Debug("recv Cmd_No_Login, [PASS]")

	case Common.Cmd_Net_Comm_Status:
		log.Debug("recv Cmd_Net_Comm_Status, [PASS]")

	case Common.Cmd_Pwd_User:
		log.Debug("recv Cmd_Pwd_User")
		ForwardToClient(&conn,&s, Common.CMDTASK)

	default:
		log.Debugf("收到%s 的自定义命令：%x",conn.RemoteAddr().String(), s.Data.Cmd)
		DESRoute(&conn,&s)
	}
}



// 目的地 路由
// 进行广播 或者 转发
// TODO:可以对链接进行是否合规检测，比如验证当前链接 是否在维护的在线列表中
func DESRoute(conn *net.Conn, s *Common.MyProtocol)  {

	des := units.BytesToDes(s.Data.Des)
	switch des {
	case 0x00:
		log.Info("Send to ComServer, [PASS]")
	case Common.Des_All:
		log.Info("Send to All")
		socket.SendToAll( units.BytesToSrc(s.Data.Src), units.BytesToCmd(s.Data.Cmd),s.Data.Data, Common.RECVTASK,0)


	case Common.Des_GB_ALL:
		log.Info("Send to GB_ALL")
		socket.SendToAll( units.BytesToSrc(s.Data.Src), units.BytesToCmd(s.Data.Cmd),s.Data.Data, Common.RECVTASK,Common.Dev_Type_GB)

	case Common.Des_UI_All:
		log.Info("Send to UI_All")
		socket.SendToAll( units.BytesToSrc(s.Data.Src), units.BytesToCmd(s.Data.Cmd),s.Data.Data, Common.RECVTASK,Common.Dev_Type_UI)


	default:
		log.Info("Send to Client：", des)


		ForwardToClient(conn,s, Common.RECVTASK)

	}
}


// 转发到客户端
// 自动判断是否在线，如果在线就转发，不在线回复客户端
func ForwardToClient(conn *net.Conn,d *Common.MyProtocol,tunnel int) {
	//判断是否在线
	client,ok :=Common.ConnList.Load(units.BytesToDes(d.Data.Des))
	if !ok{
		//未找到客户端，返回
		//目前这里只有客户端会登录，所以直接返回未找到
		sendData :=Common.Packet(10,uint16(0x0000),units.BytesToSrc(d.Data.Src), Common.Cmd_No_Login,uint32(0),nil)
		log.Errorf("ForwardToClient [ERR] %x:",sendData)
		socket.SendData(conn,sendData)
	}
	// 设备在线 发送过去
	if c,ok:=client.(*socket.Conn);ok{
		sendData :=Common.Packet(uint32(10+len(d.Data.Data)),units.BytesToSrc(d.Data.Src),units.BytesToDes(d.Data.Des),units.BytesToCmd(d.Data.Cmd),uint32(len(d.Data.Data)),d.Data.Data)

		//a :="*******************\r\n"+fmt.Sprintf("%s ForwardToClient: %x",(*conn).RemoteAddr().String(),sendData)+"\r\n"+fmt.Sprintf("SourceData:%x ",d.SourceData)

		log.Debugf("%s ForwardToClient: %x",(*conn).RemoteAddr().String(),sendData)
		//log.Printf("SourceData:%x ",d.SourceData)
		//log.Println(a)

		//单独为UI客户端进行处理一下转发的接收
		//如果转发到UI客户端的 那么直接转发改链接的cmd通道
		if c.DevType == Common.Dev_Type_UI {
			tunnel = Common.CMDTASK
		}


		if tunnel== Common.RECVTASK {
			c.RECVConn.Write(sendData)
		}else {
			c.CMDConn.Write(sendData)

		}
	}

}

