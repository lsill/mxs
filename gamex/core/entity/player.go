package entity

import (
	"mxs/api/tcp/iface"
	"mxs/proto/flat/sample/flatutil"
	"mxs/proto/flat/sample/strupro"
)

// 玩家
type Player struct {
	*Unit
	acid string
	conn iface.IConnection
}

// 创建一个玩家
func NewPlayer(conn iface.IConnection) *Player {
	unit := NewUnit()
	p := &Player{
			Unit: unit,
			acid:   "",
			conn: conn,
		}
	p.IsPlayer = true
	return p
}


func (p *Player) SendMsg(msgId uint32, data []byte) {
	p.conn.SendMsg(msgId, data)
}

// 告知客户端pid，同步已经生成的实体给客户端
func (p *Player) SyncEntity() {
	builder := flatutil.GetNewBuilder()
	posbuilder := flatutil.GetNewBuilder()
	pos := strupro.CreatePosition(posbuilder, p.X, p.Y, p.Z, p.V)
	strupro.PosMessageAddEid(builder, p.Eid)
	strupro.PosMessageAddPos(builder, pos)
}

// 附近消息（气泡消息）
func (p *Player) BubbleTalk(content []byte) {
	/*p.Entity.GetGid()
	players := p.Entity.GetCurAllPlayers()
	for i := 0; i < len(players); i++ {
		players[i].SendMsg(constcode.MsgBubbleTaik, content)
	}*/
}