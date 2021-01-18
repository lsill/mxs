package mnet

import (
	"container/list"
	"errors"
	"mxs/util/api/websocket/iface"
	"net"
	"sync"
	"sync/atomic"
)

var ErrClosed = errors.New("link.session closed")



type Session struct {
	server          iface.IServer
	id              uint64
	conn            *net.Conn
	encoder         iface.Encoder
	decoder         iface.Decoder
	closeChan       chan int
	closeFlag       int32
	closeEventMutex sync.Mutex
	metricfunc      iface.MerticFun
	State           interface{}
	closeCallbacks  *list.List
}

var globalSessionId uint64

func NewSession(conn *net.Conn, codeType iface.ICodecType, id uint64, server iface.IServer) *Session {
	return &Session{
		server: server,
		id:              atomic.AddUint64(&globalSessionId, 1),
		conn:            conn,
		encoder:         codeType.NewEncoder(*conn),
		decoder:         codeType.NewDecoder(*conn),
		closeChan:       nil,
		closeFlag:       0,
		closeEventMutex: sync.Mutex{},
		closeCallbacks:  list.New(),
		State:           nil,
	}
}

func (s *Session) Id () uint64 {
	return s.id
}

func (s *Session) Conn() *net.Conn{
	return s.conn
}

