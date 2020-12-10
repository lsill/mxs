package znet

import (
	"errors"
	"fmt"
	"net"
	"time"
	"webV/log"
	"webV/zinx/utils"
	"webV/zinx/ziface"
)

// iServer 接口实现，定义一个Server服务类

type Server struct {
	// 服务器名称
	Name string
	// tcp4 or other
	IPVersoion string
	// 服务器绑定的ip地址
	IP string
	// 服务端绑定的端口号
	Port int
	// 当前Server的消息管理模块，用来绑定msgid和对应管理方法
	msgHandler ziface.IMsgHandle
	// 当前Server的连接管理器
	ConnMgr ziface.IConnManger
	// 该Server的连接创建时hook函数
	OnConnStart func(conn ziface.IConnection)
	// 该Server的连接断开时的hook函数
	OnConnStop func(conn ziface.IConnection)
}

// ==============实现 ziface.Iserver 里的全部接口方法 =============

// 开启网络服务
func (s *Server) Start() {
	log.Release("[Start] Server listenner at IP:%v, port is %d, is starting", s.IP, s.Port)
	log.Release("[Zinx] Version:%s, MaxConn:%d, MaxPacketsize: %d", utils.GloUtil.Version, utils.GloUtil.MaxConn,
		utils.GloUtil.MaxPacketSize)
	// 0.启动worker工作池
	s.msgHandler.StarWorkerPool()
	go func() {
		// 1.获取一个TCP的Addr
		addr , err := net.ResolveTCPAddr(s.IPVersoion, fmt.Sprintf("%s:%d", s.IP,s.Port))
		if err != nil {
			log.Error("reslove tcp addr err : %v", err)
			return
		}

		//2.监听服务器地址
		listener, err := net.ListenTCP(s.IPVersoion, addr)
		if err != nil {
			log.Warn("listen %v err %v", s.IPVersoion, err)
		}
		// 已经监听成功
		log.Release("start Zinx server %v succ, now listenning...", s.Name)

		var cid uint32
		cid = 0

		// 3.启动server网络连接服务
		for {
			// 3.1 阻塞等待客户端的连接请求
			conn, err := listener.AcceptTCP()
			if err != nil {
				log.Warn("Accept err %v", err)
				continue
			}
			// 3.2 TODO Server.Start() 设置服务器最大连接控制，如果超过最大连接，那么关闭此新的连接
			if s.ConnMgr.Len() >= utils.GloUtil.MaxConn {
				conn.Close()
				log.Warn("conn is full")
				continue
			}

			// 3.3 TODO ServerStart() 处理该新链接请求的 业务 方法，此时应该有 handler 和 conn 是绑定的
			dealConn := NewConnecion(s,conn, cid, s.msgHandler)
			cid++

			// 3.4 启动当前连接的处理业务
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	log.Release("[Stop] Zinx server , name %v", s.Name)
	// TODO Server.Stop() 将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
	s.GetConnMgr().ClearConn()
}

func (s *Server) Server() {
	s.Start()

	// TODO Server.Server() 是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加


	// 阻塞， 否则主Go退出，listener的go将会退出
	for{
		time.Sleep(10 * time.Second)
	}
}

func (s *Server) AddRouter(msgid uint32, router ziface.IRouter) {
	s.msgHandler.AddRouter(msgid, router)
	log.Release("add Router succ!")
}

// 创建一个服务器句柄
func NewServer() ziface.IServer {
	s := &Server{
		Name:       utils.GloUtil.Name,
		IPVersoion: "tcp4",
		IP:         utils.GloUtil.Host,
		Port:       7777,
		msgHandler: NewMsgHandle(), // msgHandler初始化
		ConnMgr: NewConnManager(),
	}
	return s
}

// ========定义当前客户端连接的Handle apo ==========
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	//回显业务
	log.Debug("[Conn Handle] CallbackToClinet...")
	if _, err := conn.Write(data[:cnt]); err != nil {
		log.Error("write back buf err", err)
		return errors.New("CallBackToClient error")
	}
	return nil
}

// 得到连接管理
func (s *Server) GetConnMgr() ziface.IConnManger {
	return s.ConnMgr
}

// 设置该Server的连接创建时的HOOK函数
func (s *Server) SetOnConnStart(hookfunc func(ziface.IConnection)) {
	s.OnConnStart =hookfunc
}

// 设置该Server的连接断开时的hook函数
func (s *Server) SetOnConnStop(hookfunc func(connection ziface.IConnection)) {
	s.OnConnStop = hookfunc
}

// 调用连接OnConnStart hook函数
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		log.Debug("---> CallOnConnStart...")
		s.OnConnStart(conn)
	}
}

// 调用连接OnConnStop hook函数
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		log.Debug("---> CallOnConnStop")
		s.OnConnStop(conn)
	}
}