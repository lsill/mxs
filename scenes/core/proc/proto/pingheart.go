package proto

import (
	"mxs/util/api/kcp/iface"
	"mxs/util/api/kcp/mnet"
)

type PingRouter struct {
	mnet.BaseRouter
	out chan iface.IRequest
}



func (this *PingRouter) Handle(req iface.IRequest) {

}

func AddHeartBeating(typ uint32, s iface.IServer) {
	s.AddRouter(typ, &PingRouter{})
}

