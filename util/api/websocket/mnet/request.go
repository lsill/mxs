package mnet

import (
	"mxs/util/api/websocket/iface"
)

type Request struct {
	conn iface.IConnection
	pk   iface.IPackage
}

func (r *Request) GetConn() iface.IConnection {
	return r.conn
}

func (r *Request) GetPkTyp() int32 {
	return r.pk.GetTyp()
}

func (r *Request) GetData() []byte {
	return r.pk.GetData()
}

func (r *Request) SetConn(conn iface.IConnection)  {
	r.conn = conn
}



