package proto

import (
	"mxs/util/api/kcp/iface"
	"mxs/util/api/kcp/mnet"
)

type LoginOnApi struct {
	mnet.BaseRouter
}

// 发送位置信息
func (this *LoginOnApi) Handle(req iface.IRequest) {
	/*entity := strupro.GetRootAsEntity(req.GetData(), 0)
	pid, err := req.GetConnection().GetProperty("eid")
	if err != nil {
		return
	}*/
}

func AddLoginOnApi(typ uint32, s iface.IServer) {
	s.AddRouter(typ, &LoginOnApi{})
}

