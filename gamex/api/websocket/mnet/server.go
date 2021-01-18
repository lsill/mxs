package mnet

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/prometheus/common/log"
	"mxs/gamex/api/websocket/iface"
	"mxs/gamex/utils"
	logs "mxs/log"
	"net/http"
	"time"
)

type Server struct {
	name string
	ip string
	port int
	ipversion string
	msgHandler	iface.IMsgHandle
	ConnMgr iface.IConnManger
	OnConnStart func(conn iface.IConnection)
	OnConnStop func(conn iface.IConnection)
}

func NewServer() iface.IServer {
	return &Server{
		name:         "Main",
		ip:           "127.0.0.1",
		port:         2333,
		msgHandler: NewMsgHandle(),
		ConnMgr: NewConnManager(),
	}
}

func (s *Server) Server() {
	s.Start()

	// TODO Server.Server() 是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加

	for {
		time.Sleep(10 *time.Second)
	}
}
var upgrade = websocket.Upgrader{}	// 以后在配置对应参数
var Cid uint32
func (s *Server) Start() {
	logs.Release("[Start] Server listenner at addr %s:%d is startting", s.IP(), s.Port())

	go func() {
		addr := flag.String("addr", fmt.Sprintf("%s:%d", s.IP(), s.Port()), "http service address")
		http.HandleFunc("/dt", s.dtserver)
		err := http.ListenAndServe(*addr, nil)
		if err != nil {
			logs.Error("listne addr %s failed", addr)
			return
		}
	}()
}

func (s *Server)dtserver(w http.ResponseWriter, r *http.Request) {
	c, err := upgrade.Upgrade(w,r, nil)
	if err != nil {
		logs.Error("upgrade err %s", err)
		return
	}
	defer c.Close()
	if s.ConnMgr.Len() > utils.GloUtil.MaxConn {	// 此处需要修改
		c.Close()
		logs.Warn("conn is full")
		return
	}
	dealconn := NewConnecton(s, c, Cid, s.msgHandler)
	Cid++

	go dealconn.Start()
}


func(s *Server) Stop(){
	logs.Release("[Stop] server stop, name is %v", s.Name())
	s.ConnMgr.ClearConn()
}

func (s *Server) IP() string {
	return s.ip
}

func (s *Server) Port() int {
	return s.port
}

func (s *Server) IPVersion()string{
	return s.ipversion
}

func (s *Server) Name() string {
	return s.name
}

func (s *Server) AddRouter(msgid int32, router iface.IRouter) {
	s.msgHandler.AddRouter(msgid, router)
	logs.Release("add Router succ!")
}

func (s *Server) GetConnMgr() iface.IConnManger {
	return s.ConnMgr
}

// 设置该Server的连接创建时的HOOK函数
func (s *Server) SetOnConnStart(hookfunc func(iface.IConnection)) {
	s.OnConnStart =hookfunc
}

// 设置该Server的连接断开时的hook函数
func (s *Server) SetOnConnStop(hookfunc func(connection iface.IConnection)) {
	s.OnConnStop = hookfunc
}

// 调用连接OnConnStart hook函数
func (s *Server) CallOnConnStart(conn iface.IConnection) {
	if s.OnConnStart != nil {
		log.Debug("---> CallOnConnStart...")
		s.OnConnStart(conn)
	}
}

// 调用连接OnConnStop hook函数
func (s *Server) CallOnConnStop(conn iface.IConnection) {
	if s.OnConnStop != nil {
		log.Debug("---> CallOnConnStop")
		s.OnConnStop(conn)
	}
}
