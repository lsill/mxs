package proc

import (
	"mxs/scenes/core/proc/proto"
	"mxs/util/api/kcp/iface"
)

const (
	PingHeart = iota+0	// 连接心跳
)

func LoadProto(s iface.IServer) {
	proto.AddHeartBeating(PingHeart, s)
}