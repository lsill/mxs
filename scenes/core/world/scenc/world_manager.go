package scenc

import (
	"mxs/scenes/constcode"
	"mxs/scenes/core/world/aoi"
	"sync"
)

/*
	当前游戏世界的总管理模块
 */
type WorldManager struct {
	AoiMgr *aoi.AOIManager    // 当前世界地图的aoi规划管理器
	Players map[int32]*Player // 当前世界的在线玩家集合
	pLock sync.RWMutex
}

// 提供一个对外的世界管理模块句柄
var WorldMgrObj *WorldManager

// 提供WorldManager 初始化方法
func init() {
	WorldMgrObj = &WorldManager{
		AoiMgr:  aoi.NewAOIManager(constcode.AOI_MIN_X, constcode.AOI_MAX_X, constcode.AOI_CNTS_X, constcode.AOI_MIN_Y, constcode.AOI_MAX_Y, constcode.AOI_CNTS_Y),
		Players: make(map[int32]*Player),
	}
}

// 提供一个添加玩家的功能，将玩家添加进玩家信息表players
func (wm *WorldManager) AddPlayer(player *Player) {
	// 将player添加到 世界管理器中
	wm.pLock.Lock()
	wm.Players[player.Eid] = player
	wm.pLock.Unlock()

	// 将player 添加到aoi网络规划中
	wm.AoiMgr.AddToGridByPos(player.Eid, player.X,player.Y)
}

// 从玩家信息表中移除一个玩家
func (wm *WorldManager) RemovePlayerByPid(eid int32) {
	wm.pLock.Lock()
	delete(wm.Players, eid)
	wm.pLock.Unlock()
}

// 通过玩家id 获取对应玩家信息
func (wm *WorldManager) GetPlayerByEid(eid int32) *Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()
	return wm.Players[eid]
}

// 获取玩家的信息
func (wm *WorldManager) GetAllPlayers() []*Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	// 创建返回的player集合切片
	players := make([]*Player, 0, len(wm.Players))

	// 添加切片
	for _, v := range wm.Players {
		players = append(players,v)
	}
	return players
}

