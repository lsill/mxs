package aoi

import (
	"fmt"
	"mxs/api/iface"
	"mxs/mmo/core/entity"
	"sync"
)

/*
	一个地图中的格子类
 */

type Grid struct {
	GID int 	// 格子id
	MinX	int	// 格子左边界坐标
	MaxX	int // 格子右边界坐标
	MinY int	// 格子上边界坐标
	MaxY	int // 格子下边界坐标
	entityIDs map[int]*entity.Entity	// 当前格子内实体成员
	players map[int]*entity.Player	// 当前格子内的玩家
	units map[int]*entity.Unit	// 当前格子内的实体
	entityLock sync.RWMutex	// 实体map锁
	unitLock sync.RWMutex	// 单位锁
	playLock sync.RWMutex	// 玩家锁
}

// 初始化一个格子
func NewGrid(gId, minX, maxX, minY,maxY int) *Grid {
	return &Grid{
		GID:        gId,
		MinX:       minX,
		MaxX:       maxX,
		MinY:       minY,
		MaxY:       maxY,
		entityIDs:  make(map[int]*entity.Entity),
	}
}

// 向当前格子中添加一个实体
func (g *Grid) AddEntity(entityid int) {
	g.entityLock.Lock()
	defer g.entityLock.Unlock()
	g.entityIDs[entityid] = entity.NewEntity()
}

// 从格子中删除一个实体
func (g *Grid) RemoveEntity(entityid int) {
	g.entityLock.Lock()
	defer g.entityLock.Unlock()
	delete(g.entityIDs, entityid)
}


// 向当前格子中添加一个单位
func (g *Grid) AddUnit(entityid int) {
	g.unitLock.Lock()
	defer g.unitLock.Unlock()
	g.units[entityid] = entity.NewUnit()
}

// 从格子中删除一个单位
func (g *Grid) RemoveUnit(entityid int) {
	g.unitLock.Lock()
	defer g.unitLock.Unlock()
	delete(g.units, entityid)
}

// 向当前格子中添加一个角色
func (g *Grid) AddPlayer(entityid int, conn iface.IConnection) {
	g.playLock.Lock()
	defer g.playLock.Unlock()
	g.players[entityid] = entity.NewPlayer(conn)
}

// 从格子中删除一个角色
func (g *Grid) RemovePlayer(entityid int) {
	g.playLock.Lock()
	defer g.playLock.Unlock()
	delete(g.players, entityid)
}


// 得到当前格子的所有实体id
func (g *Grid) GetAllEntityIDs() []int {
	g.entityLock.RLock()
	defer g.entityLock.RUnlock()
	entityIDs := make([]int, 0 , len(g.entityIDs))
	for k, _ := range g.entityIDs {
		entityIDs = append(entityIDs, k)
	}
	return entityIDs
}
// 得到当前格子内的所有实体
func (g *Grid) GetAllEntitys() []*entity.Entity {
	g.entityLock.RLock()
	defer g.entityLock.RUnlock()
	entitys := make([]*entity.Entity, 0 , len(g.entityIDs))
	for _, v := range g.entityIDs {
		entitys = append(entitys, v)
	}
	return entitys
}

// 得到当前格子的所有单位id
func (g *Grid) GetAllUnitIDs() []int {
	g.unitLock.RLock()
	defer g.unitLock.RUnlock()
	unitIds := make([]int, 0 , len(g.units))
	for k, _ := range g.units {
		unitIds = append(unitIds, k)
	}
	return unitIds
}
// 得到当前格子内的所有单位
func (g *Grid) GetAllUnits() []*entity.Unit {
	g.unitLock.RLock()
	defer g.unitLock.RUnlock()
	units := make([]*entity.Unit, 0 , len(g.units))
	for _, v := range g.units {
		units = append(units, v)
	}
	return units
}

// 得到当前格子的所有角色id
func (g *Grid) GetAllPlayerIds() []int {
	g.playLock.RLock()
	defer g.playLock.RUnlock()
	playerids := make([]int, 0 , len(g.players))
	for k, _ := range g.entityIDs {
		playerids = append(playerids, k)
	}
	return playerids
}
// 得到当前格子内的所有角色
func (g *Grid) GetAllPlayers() []*entity.Player {
	g.playLock.RLock()
	defer g.playLock.RUnlock()
	players := make([]*entity.Player, 0 , len(g.entityIDs))
	for _, v := range g.players {
		players = append(players, v)
	}
	return players
}


// 重写打印信息方法
func (g *Grid) String() string {
	return fmt.Sprintf("Gird id :%d, minX:%d, maxX:%d, minY:%d, maxY:%d, entityIDs:%v", g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.entityIDs)
}
