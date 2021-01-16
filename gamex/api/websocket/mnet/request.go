package mnet

import (
	"mxs/gamex/api/websocket/iface"
)

type Request struct {
	conn iface.IConnection
	pk      iface.IPackage
}

func (r *Request) GetConn() iface.IConnection {
	return r.conn
}

func (r *Request) GetPkTyp() int32 {
	return r.GetPkTyp()
}

func (r *Request) GetData() []byte {
	return r.GetData()
}