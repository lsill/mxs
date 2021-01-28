package aoi

import (
	"fmt"
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
	entityIDs map[int32]bool	// 当前格子内实体成员
	entityLock sync.RWMutex	// 实体map锁
}

// 初始化一个格子
func NewGrid(gId, minX, maxX, minY,maxY int) *Grid {
	return &Grid{
		GID:        gId,
		MinX:       minX,
		MaxX:       maxX,
		MinY:       minY,
		MaxY:       maxY,
		entityIDs:  make(map[int32]bool),
	}
}

// 向当前格子中添加一个实体
func (g *Grid) AddEntity(entityid int32) {
	g.entityLock.Lock()
	defer g.entityLock.Unlock()
	g.entityIDs[entityid] = true
}

// 从格子中删除一个实体
func (g *Grid) RemoveEntity(entityid int32) {
	g.entityLock.Lock()
	defer g.entityLock.Unlock()
	delete(g.entityIDs, entityid)
}



// 得到当前格子的所有实体id
func (g *Grid) GetAllEntityIDs() []int32 {
	g.entityLock.RLock()
	defer g.entityLock.RUnlock()
	entityIDs := make([]int32, 0 , len(g.entityIDs))
	for k, _ := range g.entityIDs {
		entityIDs = append(entityIDs, k)
	}
	return entityIDs
}

// 重写打印信息方法
func (g *Grid) String() string {
	return fmt.Sprintf("Gird id :%d, minX:%d, maxX:%d, minY:%d, maxY:%d, entityIDs:%v", g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.entityIDs)
}

