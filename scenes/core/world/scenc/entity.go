package scenc

import (
	"math/rand"
	"mxs/log"
	"mxs/scenes/constcode"
	"mxs/scenes/proto/flat/flatbuffers"
	"mxs/util/api/kcp/iface"
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
	IsPlayer bool // 是否玩家
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
/*	builder := flatbuffers.NewBuilder(1024)
	strupro.EntityAddEid(builder, id)
	pos := strupro.CreatePosition(builder,1.0,1.0,1.0)
	strupro.EntityAddPos(builder, pos)
	entity := strupro.EntityEnd(builder)*/
	EidGen++
	IdLock.Unlock()
	en := &Entity{
		Eid: id,
		X:   float32(50 + rand.Intn(10)),
		Y:   float32(100 + rand.Intn(17)),
		Z:   0,
		V:   0,
	}
	return en
}

func (en *Player) SyncPlayers() {
	pids := WorldMgrObj.AoiMgr.GetEIDsByPos(en.X,en.Y)
	log.Debug("len pids is %d", len(pids))
	players := make([]*Player, 0, len(pids))
	for _,pid := range pids {
		player := WorldMgrObj.GetPlayerByEid(pid)
		if player.acid != "" {
			players = append(players, player)
		}
	}
	builder := flatbuffers.NewBuilder(20480)
	builder.Finish(GenEntityProto(builder, en.Entity))
	bytes := builder.Bytes[builder.Head():]
	players = append(players, en)
	for _,player := range players {
		player.SendMsg(constcode.PositionOther, bytes)
	}
	builders := flatbuffers.NewBuilder(20480)
	builders.Finish(GenPlayersProto(builders, players))
	bytes = builders.Bytes[builders.Head():]
	en.SendMsg(constcode.PositionMine, bytes)
}

/*
builder := flatbuffers.NewBuilder(2000)
h := flatutil.NewFlatBufferHelper(builder, 32)
id := h.Pre(builder.CreateString("hello kcp"))
strupro.TestMessageStart(builder)
strupro.TestMessageAddTeststr(builder, h.Get(id))
strupro.TestMessageEnd(builder)
dp := mnet.NewDataPack()
bytes :=  builder.Bytes[builder.Head():]*/

// 当客户端建立链接的时候的hook函数
func OnConnectionAdd(conn iface.IKConnection) {
	// 创建一个玩家
	player := NewPlayer(conn)
	// 将当前新上线玩家添加到worldmanager中
	WorldMgrObj.AddPlayer(player)
	conn.SetProperty("eid", player.Eid)
	// 同步当前的玩家id坐标 给客户端
	player.SyncPlayers()
	log.Release("====> player eid = %d", player.Eid)
}

func OnConnectionDel(conn iface.IKConnection) {
	obj, _ := conn.GetProperty("eid")
	WorldMgrObj.RemovePlayerByPid(obj.(int32))
}




