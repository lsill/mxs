package entity

import (
	"math/rand"
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







