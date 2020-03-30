package business

import "C"
import (
	"errors"
	"fmt"
	"github.com/Code-Fight/golog"
	"math"
	"net"
	"socketserver/Common"
	"socketserver/socket"
	"socketserver/units"
	"strings"
	"sync"
)

// 设备注册
func Reg(conn net.Conn, s *Common.MyProtocol,closeChannel chan struct{}) {

	if len(s.Data.Data)!=5{
		log.Error("设备注册过程中，缺少数据")
		return
	}

	ZBM:=units.BytesToString (s.Data.Data[2:])

	if units.BytesToSrc(s.Data.Src) == 0xffff {
		c := socket.Conn{}
		c.ZBM = ZBM

		//设备类型 从data中来
		c.DevType = units.BytesToSrc(s.Data.Data[:2])

		//分配设备号 并创建socket Conn 添加到Clients中
		_,disOk := distributionID(&c,&conn)
		if disOk!=nil{
			if conn!=nil{
				conn.Close()
			}
			log.Error("设备分配过程中失败:"+disOk.Error())
			return
		}


		//回复client ID
		c.ReplyDevId(c.RECVConn,c.DevId)

		Common.ConnListIp.LoadOrStore(conn.RemoteAddr().String(),ZBM)

		//给当前连接器一个心跳线程
		go HeardEvent(c.RECVConn,c.DevId,closeChannel)


		return
	} else if tempRoom, ok := Common.ClientList.Load(ZBM); ok {

		//首先判断是否存在房间号，然后判断改房间是否存在已经注册的RECV链接

		room,_:=tempRoom.(sync.Map)

		if roomClient,roomOk:=room.Load(units.BytesToSrc(s.Data.Src));roomOk{

			// 自带设备号过来的连接  添加到 发送数据CMD通道
			c, ok := roomClient.(*socket.Conn)

			if ok {
				c.CMDConn = conn
				//回复client ID
				c.ReplyDevId(c.CMDConn,c.DevId)

				Common.ConnListIp.LoadOrStore(conn.RemoteAddr().String(),ZBM)

				//通知所有设备有设备上线
				statusBytes := GenDevStatusBytes(uint32(c.DevId),uint32(c.DevType),1,strings.Split(conn.RemoteAddr().String(),":")[0])
				socket.SendToRoom(0x0000, Common.Cmd_Net_Comm_Status,statusBytes, Common.RECVTASK,0,ZBM)

				//通知新上线的设备 存在哪些已在线的设备
				GetAllOnlineDev(c.DevId,ZBM)

				//开启心跳
				go HeardEvent(c.CMDConn,c.DevId,closeChannel)



			} else {
				if conn!=nil{
					conn.Close()
				}
				log.Error("设备注册过程中，未找到发送数据相关的接收通道")
				return

			}

		}else {
			log.Error("CMD通道注册失败")
			if conn!=nil{
				conn.Close()
			}
			return
		}


	} else {
		// 没有带设备号 并 没有设置来源0xffff的 直接关闭
		log.Error("设备没有申请设备id，以及设备未注册到clients中")
		if conn!=nil{
			conn.Close()
		}
		return
	}

}

// 获取所有在线设备
func GetAllOnlineDev(CurrClientID uint16,ZBM string)  {


	Room,roomOk:=Common.ClientList.Load(ZBM)
	if  roomOk {
		RoomClient,RoomClientOk:=Room.(sync.Map)

		if !RoomClientOk{
			return
		}



		_c,ok:=RoomClient.Load(CurrClientID)
		if ok {
			CurrClient,succ:=_c.(*socket.Conn)

			if succ && CurrClient!=nil{

				RoomClient.Range(func(key, value interface{}) bool {


					client,ok :=value.(*socket.Conn)

					if ok{
						if client.DevId !=CurrClientID{
							//_,_ :=key.(uint16)

							data:=GenDevStatusBytes(uint32(client.DevId),uint32(client.DevType),1,strings.Split(client.RECVConn.RemoteAddr().String(),":")[0])

							sendData := Common.Packet(uint32(len(data)+10),0x0000,CurrClientID, Common.Cmd_Net_Comm_Status,uint32(len(data)),data)
							log.Debugf("GetAllDevs:%x",sendData)
							if CurrClient.RECVConn!=nil{
								_,ok :=CurrClient.RECVConn.Write(sendData)
								if ok!=nil {
									log.Error("发送数据失败：",ok.Error())
								}
							}

						}


					}


					return true
				})
			}

		}
	}




}

// status 1=on 0=off
func GenDevStatusBytes(id uint32, dtype uint32, status uint32,ip string) []byte {
	n := Common.NETCOMM_STATUS{}
	n.Init(id,dtype,status,ip)
	return 	n.ToBytes()

}


//分配ID
//此处加锁 防止并发时 多个设备一个id号
func distributionID(conn *socket.Conn,netConn *net.Conn) (n uint16, err error) {
	Common.RegMutex.Lock()
	defer Common.RegMutex.Unlock()


	//取站报码
	zbmClients,zbmOK:=Common.ClientList.Load(conn.ZBM)

	//设备ID跟设备类型相关
	switch conn.DevType {
	case Common.Dev_Type_DX, Common.Dev_Type_GB, Common.Dev_Type_YW, Common.Dev_Type_DB:


		//如果找到相应的房间
		if zbmOK{
			zbmCliecnts,_:=zbmClients.(sync.Map)
			_,ok :=zbmCliecnts.Load(conn.DevType)
			if ok{
				return 0,errors.New(fmt.Sprintf("%x,该设备已经在线，不能重复登录.",conn.DevType))
			}
			conn.DevId = conn.DevType
		}else{
			//找不到房间，创建房间 并分配房间号
			conn.DevId = conn.DevType
		}


	case Common.Dev_Type_UI:
		//0x10000


		//如果找到相应的房间
		if zbmOK{

			zbmCliecnts,_:=zbmClients.(sync.Map)
			//roomClient,ok :=zbmCliecnts.Load(conn.DevType)
			clientId := uint16(4096)
			for{
				_,ok:= zbmCliecnts.Load(clientId)
				if !ok{
					conn.DevId = clientId
					break
				}
				if clientId > math.MaxUint16{
					return 0,errors.New(fmt.Sprintf("%x,该类型设备已到达上线.",conn.DevType))
				}
				clientId++
			}
		}





	default:
		return 0, errors.New("未定义设备类型")
	}



	// 设置链接  保存到clients中
	conn.RECVConn = *netConn
	conn.CMDConn = nil

	//保存链接
	if zbmOK{
		//存在房间，保存到房间中
		zbmCliecnts,_:=zbmClients.(sync.Map)
		zbmCliecnts.Store(conn.DevId,conn)

	}else {
		//房间不存在 创建房间 并保存
		room :=sync.Map{}
		room.Store(conn.DevId,conn)
		Common.ClientList.Store(conn.ZBM,room)
	}


	return conn.DevId,nil
}
