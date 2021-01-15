package mnet

import (
	"flag"
	"fmt"
	logs "mxs/log"
	"net/http"
	"time"
)

const maxSessionId = 1000

type Server struct {
	name string
	ip string
	port int
	ipversion string
 	//codecType iface.ICodecType

	// about sessions
	maxSessionId uint64
}

func NewServer() *Server {
	return &Server{
		name:         "Main",
		ip:           "127.0.0.1",
		port:         2333,
		maxSessionId: maxSessionId,
	}
}

func (s *Server) Server() {
	s.Start()

	// TODO Server.Server() 是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加

	for {
		time.Sleep(10 *time.Second)
	}
}

func (s *Server) Start() {
	logs.Release("[Start] Server listenner at addr %s:%d is startting", s.IP(), s.Port())

	go func() {
		addr := flag.String("addr", fmt.Sprintf("%s:%d", s.IP(), s.Port()), "http service address")
		err := http.ListenAndServe(*addr, nil)
		if err != nil {
			logs.Error("listne addr %s failed", addr)
			return
		}
	}()
}

func(s *Server) Stop(){

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

