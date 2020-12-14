package entity

import (
	"math/rand"
	"mxs/mmo/core/aoi"
	"mxs/mmo/core/worldmanager"
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


// 得到当前实体所在格子的id
func (en *Entity) GetGid() int{
	return worldmanager.WorldMgrObj.AoiMgr.GetGIDByPos(en.X, en.Y)
}

// 得到当当前实体附近所有格子
func (en *Entity) GetCurAllGirds() []*aoi.Grid {
	return worldmanager.WorldMgrObj.AoiMgr.GetSurroundGridsByGid(en.GetGid())
}

// 得到当前实体附近所有实体
func (en *Entity) GetCurAllEntitys() []*Entity {
  	entitys := make([]*Entity, 0 , 1000)
  	girds := en.GetCurAllGirds()
  	for i := 0; i < len(girds); i++{
  		entitys = append(entitys, girds[i].GetAllEntitys()...)
	}
	return entitys
}

// 得到当前角色附近所有玩家
func (en *Entity) GetCurAllPlayers() []*Player {
	entitys := make([]*Player, 0 , 1000)
	girds := en.GetCurAllGirds()
	for i := 0; i < len(girds); i++{
		entitys = append(entitys, girds[i].GetAllPlayers()...)
	}
	return entitys
}




