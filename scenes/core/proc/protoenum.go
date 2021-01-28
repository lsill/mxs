package proc

import (
	"mxs/scenes/constcode"
	"mxs/scenes/core/proc/proto"
	"mxs/scenes/core/world/scenc"
	"mxs/util/api/kcp/iface"
)



func LoadProto(s iface.IServer) {
	s.SetOnConnStart(scenc.OnConnectionAdd)
	s.SetOnConnStop(scenc.OnConnectionDel)


	proto.AddHeartBeating(constcode.PingHeart, s)
	proto.AddMoveApi(constcode.PositionMine, s)
}

