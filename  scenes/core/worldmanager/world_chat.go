package worldmanager

import (
	"mxs/gamex/proto/flat/sample/strupro"
	"mxs/log"
	"mxs/util/api/tcp/iface"
	"mxs/util/api/tcp/mnet"
)

// 世界聊天 路由业务
type WorldChat struct {
	mnet.BaseRouter
}

func (*WorldManager) Handler(request iface.IRequest) {
	msg := strupro.GetSizePrefixedRootAsTalkMessage(request.GetData(), 32)

	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		log.Error("GetProperty pid error %v", pid)
		request.GetConnection().Stop()
		return
	}

	player := WorldMgrObj.GetPlayerByEid(pid.(int32))

	player.BubbleTalk(msg.Table().Bytes)
}

