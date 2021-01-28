package aoi

import (
	"fmt"
	"mxs/log"
)


/*
	AOI管理模块
 */
type AOIManager struct {
	MinX int            // 区域左边界坐标
	MaxX int            // 区域右边界坐标
	CntsX int           // x方向格子的数量
	MinY int            // 区域上边界坐标
	MaxY int            // 区域下边界坐标
	GntsY int           // y方向的格子数量
	grids map[int]*Grid // 当前区域都有那些格子，key：格子id，value：格子对象
}

/*
	初始化一个AOI区域
 */
func NewAOIManager(minX, maxX, cntsX, minY, maxY, cntsY int) *AOIManager {
	aoiMgr := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		CntsX: cntsX,
		MinY:  minY,
		MaxY:  maxY,
		GntsY:  cntsY,
		grids: make(map[int]*Grid),
	}
	// 给AOI初始化区域中所有的格子
	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			// 计算格子id
			// 格子编号：id = idy * nx + idx (利用格子坐标得到格子编号)
			gid := y * cntsX + x

			// 初始化一个格子放在AOI中的map里，key是当前格子的id
			aoiMgr.grids[gid] = NewGrid(gid,
				aoiMgr.MinX + x * aoiMgr.gridWidth(),
				aoiMgr.MinX + (x+1) * aoiMgr.gridWidth(),
				aoiMgr.MinY + y * aoiMgr.gridLength(),
				aoiMgr.MinY + (y+1)*aoiMgr.gridLength())
		}
	}
	return aoiMgr
}

// 得到每个格子在x轴方向的宽度
func (m *AOIManager) gridWidth() int {
	return (m.MaxX - m.MinX) / m.CntsX
}

// 得到每个格子在y轴方向的长度
func (m *AOIManager) gridLength() int {
	return (m.MaxY - m.MinY) / m.GntsY
}

// 打印信息方法
func (m *AOIManager) String() string {
	s:=fmt.Sprintf("AOIManager:\nminX:%d, maxX:%d, cntsX:%d, minY:%d, maxY:%d, cntsY:%d\n", m.MinX, m.MaxX,m.CntsX,m.MinY,m.MaxY,m.GntsY)
	for _, grid := range m.grids {
		s += fmt.Sprintln(grid)
	}
	return s
}

// 根据格子gid得到当前周边的九宫格信息
func (m *AOIManager) GetSurroundGridsByGid(gID int) []*Grid {
	grids := make([]*Grid, 0, len(m.grids))
	// 判断gid是否存在
	if _, ok := m.grids[gID]; !ok {
		return grids
	}
	// 将当前gid添加到九宫格中
	grids = append(grids, m.grids[gID])

	// 根据gid得到当前格子所在x轴编号
	idx := gID % m.CntsX
	// 判断当前idx左边是否还有格子
	if idx > 0 {
		grids = append(grids, m.grids[gID -1])
	}
	// 判断当前idx右边是否还有格子
	if idx < m.CntsX - 1 {
		grids = append(grids, m.grids[gID + 1])
	}
	// 将x轴当前的格子都取出，进行遍历，在分别得到每个格子的上下是否还有格子

	// 得到当前x轴的格子id集合
	gidsX := make([]int, 0, len(grids))
	for _,v := range grids {
		gidsX = append(gidsX, v.GID)
	}

	// 遍历x轴格子
	for _, v := range gidsX {
		// 计算该格子处于第几列
		idy := v / m.CntsX
		// 判断当前的idy上边是否还有格子
		if idy > 0 {
			grids = append(grids, m.grids[v - m.CntsX])
		}
		// 判断当前的idy下边是否还有格子
		if idy < m.GntsY -1 {
			grids = append(grids, m.grids[v + m.CntsX])
		}
	}
	return grids
}

// 通过横纵坐标获取对应的格子id
func (m *AOIManager) GetGIDByPos(x, y float32) int {
	gx := (int(x) - m.MinX) / m.gridWidth()
	gy := (int(y) - m.MinY) / m.gridLength()
	return gy * m.CntsX + gx
}

// 通过横纵坐标获取周边九宫格内的全部entityids

func (m *AOIManager) GetEIDsByPos(x, y float32) []int32 {
	// 根据横纵坐标获取当前坐标属于哪个格子
	gid := m.GetGIDByPos(x, y)

	// 根据格子id的得到周边九宫格的信息
	grids := m.GetSurroundGridsByGid(gid)
	entityids := make([]int32, 0, len(grids))
	for _, v := range grids {
		entityids = append(entityids, v.GetAllEntityIDs()...)
		log.Debug("====> grid ID : %d, pids: %v ====", v.GID, v.GetAllEntityIDs())
	}
	return entityids
}

// 通过GID获取当前格子的全部entituid
func (m *AOIManager) GetEidsByGid(gid int) (entityIds []int32) {
	entityIds = m.grids[gid].GetAllEntityIDs()
	return
}

// 移除一个格子中的entityid
func (m *AOIManager) RemoveEidFromFGrid(eid int32, gid int) {
	m.grids[gid].RemoveEntity(eid)
}

// 添加一个entity到一个格子中
func(m *AOIManager) AddEidToGrid(eid int32, gid int) {
	m.grids[gid].AddEntity(eid)
}

// 通过横纵坐标添加一个entity到一个格子中
func (m *AOIManager) AddToGridByPos(eid int32, x, y float32) {
	gid := m.GetGIDByPos(x, y)
	grid := m.grids[gid]
	grid.AddEntity(eid)
}

// 通过横纵坐标从格子中删除一个entity
func (m *AOIManager) RemoveFromGridByPos(eid int32, x, y float32) {
	gid := m.GetGIDByPos(x, y)
	grid := m.grids[gid]
	grid.RemoveEntity(eid)
}
