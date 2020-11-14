package znet

import (
	"fmt"
	"net"
	"time"
	"zinx/ziface"
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
}

// ==============实现 ziface.Iserver 里的全部接口方法 =============

// 开启网络服务
func (s *Server) Start() {
	fmt.Printf("[Start] Server listenner at IP:%s, port is %d, is starting\n", s.IP, s.Port)
	go func() {
		// 1.获取一个TCP的Addr
		addr , err := net.ResolveTCPAddr(s.IPVersoion, fmt.Sprintf("%s:%d", s.IP,s.Port))
		if err != nil {
			fmt.Println("reslove tcp addr err : ", err)
			return
		}

		//2.监听服务器地址
		listener, err := net.ListenTCP(s.IPVersoion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersoion, " err", err)
		}
		// 已经监听成功
		fmt.Println("start Zinx server", s.Name, "succ, now listenning...")
		// 3.启动server网络连接服务
		for {
			// 3.1 阻塞等待客户端的连接请求
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			// 3.2 TODO Server.Start() 设置服务器最大连接控制，如果超过最大连接，那么关闭此新的连接

			// 3.3 TODO ServerStart() 处理该新链接请求的 业务 方法，此时应该有 handler 和 conn 是绑定的

			//这里暂时做一个最大512字节的回显服务
			go func () {
				//不断的循环从客户端获取数据
				for  {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Println("recv buf err ", err)
						return
					}
					//回显
					if _, err := conn.Write(buf[:cnt]); err !=nil {
						fmt.Println("write back buf err ", err)
						continue
					}
				}
			}()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[Stop] Zinx server , name", s.Name)

	// TODO Server.Stop() 将其他需要清理的连接信息或者其他信息 也要一并停止或者清理

}

func (s *Server) Server() {
	s.Start()

	// TODO Server.Server() 是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加


	// 阻塞， 否则主Go退出，listener的go将会退出
	for{
		time.Sleep(10 * time.Second)
	}
}

// 创建一个服务器句柄

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:       name,
		IPVersoion: "tcp4",
		IP:         "0.0.0.0",
		Port:       7777,
	}
	return  s
}