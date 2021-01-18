package iface

import (
	"github.com/gorilla/websocket"
	"net"
)

type IConnection interface {
	Start()
	Stop()
	Conn() *websocket.Conn
	ConnId() uint32
	RemoteAddr()  net.Addr
	SendMsg(msgId uint32, data[]byte) error
}

