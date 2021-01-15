package mnet

import (
	"mxs/gamex/api/websocket/iface"
)

type Request struct {
	session *iface.ISession
	pk      iface.IPacket
}

func (r *Request) Getsession() *iface.ISession {
	return r.session
}

func (r *Request) GetPkTyp() int32 {
	return r.pk.Typ()
}

func (r *Request) GetData() []byte {
	return r.pk.Data()
}