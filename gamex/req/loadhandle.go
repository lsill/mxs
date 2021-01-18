package req

import "mxs/util/api/websocket/iface"

func Load(s iface.IServer) {
	s.AddRouter(0, &HelloRouter{})
}
