package mnet

import (
	"errors"
	"mxs/api/tcp/iface"
	"mxs/log"
	"sync"
)

/*
	链接管理模块
 */
type ConnManager struct {
	connections map[uint32]iface.IConnection //管理的链接信息
	connLock	sync.RWMutex                 // 读写链接的读写锁
}

/*
	创建一个链接管理
 */
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]iface.IConnection),
	}
}

// 添加链接
func (cmg *ConnManager) Add(conn iface.IConnection) {
	// 保护共享资源map，加写锁
	cmg.connLock.Lock()
	defer cmg.connLock.Unlock()

	// 将conn链接添加到ConnMananger中
	cmg.connections[conn.GetConnID()] = conn

	log.Debug("connection add conn %d to ConnManager success, conn num is %d", conn.GetConnID(),cmg.Len())
}

// 删除链接
func (cmg *ConnManager) Remove(conn iface.IConnection) {
	cmg.connLock.Lock()
	defer cmg.connLock.Unlock()
	// 删除链接信息
	delete(cmg.connections, conn.GetConnID())
	log.Debug("connection Remove ConnID %d success, conn num is %d",conn.GetConnID(), cmg.Len())
}

// 获取当前连接数量
func (cmg *ConnManager) Len() int {
	return len(cmg.connections)
}

// 通过ConnID获取链接
func (cmg *ConnManager) Get(connid uint32) (iface.IConnection,error){
	// 保护共享资源map，加读锁
	cmg.connLock.Lock()
	defer cmg.connLock.Unlock()

	if conn, ok := cmg.connections[connid]; ok {
		return conn, nil
	}
	return nil, errors.New("connection not found")
}

// 停止并清除所有链接
func (cmg *ConnManager) ClearConn() {
	// 保护共享资源map，加写锁
	cmg.connLock.Lock()
	defer cmg.connLock.Unlock()

	// 停止并删除全部的连接信息
	for connID, conn := range cmg.connections{
		// 停止
		conn.Stop()
		delete(cmg.connections,connID)
	}
	log.Release("Clear ALL Connections successful, conn num is %d", cmg.Len())
}