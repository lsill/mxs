package entity

import (
	"math/rand"
	"mxs/api/iface"
	"mxs/mmo/proto/flat/flatbuffers"
	"mxs/mmo/proto/flat/sample/flatutil"
	"mxs/mmo/proto/flat/sample/strupro"
	"sync"
)

// 实体
type Entity struct {
	Eid	int32	// 实体id
	X float32	// 平面x坐标
	Y float32	// 平面y坐标
	Z float32	// 高度
	V float32	// 旋转角度
	W int32		// 重量
}

// 玩家
type Player struct {
	Entity
	acid string
	conn iface.IConnection
}

/*
	entity 实体id生成器
 */
var EidGen int32 = 1	// 用来生成实体id的计数器
var IdLock sync.Mutex	// 保护实体id唯一的互斥锁

// 创建一个实体
func NewEntity() *Entity {
	// 生成一个实体id
	IdLock.Lock()
	id := EidGen
	EidGen++
	IdLock.Unlock()
	en := &Entity{
		Eid: id,
		X:   float32(160 + rand.Intn(10)),
		Y:   float32(134 + rand.Intn(17)),
		Z:   0,
		V:   0,
	}
	return en
}

// 创建一个玩家
func NewPlayer(conn iface.IConnection) *Player {
	IdLock.Lock()
	id := EidGen
	EidGen++
	IdLock.Unlock()
	p := &Player{
		Entity: Entity{
			Eid: id,
			X:   float32(160 + rand.Intn(10)),
			Y:   float32(134 + rand.Intn(17)),
			Z:   0,
			V:   0,
		},
		acid:   "",
		conn: conn,
	}
	return p
}

func (p *Player) SendMsg(msgId uint32, data flatbuffers.FlatBuffer) {
	p.conn.SendMsg(msgId, data.Table().Bytes)
}

// 告知客户端pid，同步已经生成的实体给客户端
func (p *Player) SyncEntity() {
	builder := flatutil.GetNewBuilder()
	posbuilder := flatutil.GetNewBuilder()
	pos := strupro.CreatePosition(posbuilder, p.X, p.Y, p.Z, p.V)
	strupro.PosMessageAddEid(builder, p.Eid)
	strupro.PosMessageAddPos(builder, pos)
}