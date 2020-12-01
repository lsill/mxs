package main

import (
	"zinx/log"
	"zinx/ziface"
	"zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter	// 一定要先基础BaseRouter
}

// Test PreHandle
/*func (this *PingRouter) PreHandle(request ziface.IRequest){
	log.Debug("Call Router PreHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping....\n"))
	if err != nil {
		log.Error("call back ping ping ping error")
	}
}*/

// TestHandle
func (this *PingRouter) Handle(request ziface.IRequest) {
	log.Debug("Call PingRouter Handle")
	log.Debug("recv from client,msgid:%d, data=%s", request.GetMsgID(), string(request.GetData()))
	err := request.GetConnection().SendMsg(0, []byte("ping...ping...ping..."))
	if err != nil {
		log.Error("sendmsg err %v", err)
	}
}

// TestPostHandle
/*func (this *PingRouter) PostHandle(request ziface.IRequest) {
	log.Debug("Call Router PostHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("After ping....\n"))
	if err != nil {
		log.Error("Call back ping ping ping error")
	}
}*/

type HelloZinxRouter struct {
	znet.BaseRouter
}

func (this *HelloZinxRouter) Handle(request ziface.IRequest) {
	log.Debug("Call HelloZinxRouter Handle")
	log.Debug("recv from client,msgid:%d, data=%s", request.GetMsgID(), string(request.GetData()))
	err := request.GetConnection().SendMsg(1, []byte("hello Zinx Router 0.6"))
	if err != nil {
		log.Error("err is %v", err)
	}
}

func DoConnectionBegin(conn ziface.IConnection) {
	log.Debug("DoConnectionBegin is Called...")
	conn.SetProperty("Name", "Lsill")
	conn.SetProperty("Home", "https://www.baidu.com/")
	err := conn.SendMsg(2, []byte("DoConnection Begin..."))
	if err != nil {
		log.Error("err is %v", err)
	}
}

func DoConnectionLost(conn ziface.IConnection) {
	log.Debug("DoConnectionLost is Called...")
	if name, err := conn.GetProperty("Name"); err == nil {
		log.Debug("Conn property name is %s", name.(string))
	} else {
		log.Error("GetProperty name err is %v", err)
	}
	if home, err := conn.GetProperty("Home"); err == nil {
		log.Debug("Conn property Home is %v", home.(string))
	} else {
		log.Error("GetProperty home err is %v", err)
	}
	log.Debug("DoConnectionLost is Called...")
}

func main() {
	s := znet.NewServer()
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})
	s.Server()
}

