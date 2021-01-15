package iface

import "net"

type ISession interface {
	Id() uint64
	Conn() net.Conn
	IsClosed() bool
	Close() bool
	Receive(msg interface{}) error
	Send(msg interface{}) error
}
