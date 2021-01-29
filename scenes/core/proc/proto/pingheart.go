package proto

import (
	logs "mxs/log"
	"mxs/util/api/kcp/iface"
	"mxs/util/api/kcp/mnet"
)

type PingRouter struct {
	mnet.BaseRouter
}


// 心跳连接不实现
func (this *PingRouter) Handle(req iface.IRequest){
	logs.Debug("ping ping ping")
}

func AddHeartBeating(typ uint32, s iface.IServer) {
	s.AddRouter(typ, &PingRouter{})
}

