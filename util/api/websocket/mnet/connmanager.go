package mnet

import (
	"errors"
	logs "mxs/log"
	"mxs/util/api/websocket/iface"
	"sync"
)

type ConnManager struct {
	connections map[uint32]iface.IConnection
	connLock sync.RWMutex	// 读写锁
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]iface.IConnection),
		connLock:    sync.RWMutex{},
	}
}

func (cmg *ConnManager) Add(conn iface.IConnection) {
	cmg.connLock.Lock()
	defer cmg.connLock.Unlock()
	cmg.connections[conn.ConnId()]=conn
	logs.Debug("connection add conn %d to connManager success, conn num is %d", conn.ConnId(), cmg.Len())
}

func (cmg *ConnManager) Len() int{
	return len(cmg.connections)
}

func (cmg *ConnManager) Remove(conn iface.IConnection) {
	cmg.connLock.Lock()
	defer cmg.connLock.Unlock()
	delete(cmg.connections, conn.ConnId())
	logs.Debug("connmamager remove connid %d, success, conn num is %d", conn.ConnId(), cmg.Len())
}

func (cmg *ConnManager) Get(connid uint32) (iface.IConnection, error) {
	cmg.connLock.RLock()
	defer cmg.connLock.RUnlock()
	if conn, ok := cmg.connections[connid]; ok {
		return conn, nil
	}
	return nil, errors.New("connection not found")
}

func (cmg *ConnManager) ClearConn() {
	cmg.connLock.Lock()
	defer cmg.connLock.Unlock()
	for connid, conn := range cmg.connections{
		conn.Stop()
		delete(cmg.connections, connid)
	}
	logs.Release("Clear All connections successful!")
}