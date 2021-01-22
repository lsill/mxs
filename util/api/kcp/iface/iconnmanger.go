package iface

/*
	连接管理抽象层
 */

type IConnManger interface {
	Add(conn IKConnection)                  // 添加链接
	Remove(conn IKConnection)               // 删除链接
	Get(connID uint32)(IKConnection, error) // 利用ConnID获取链接
	Len() int                              // 获取当前连接
	ClearConn()                            // 删除并停止所有链接
}


