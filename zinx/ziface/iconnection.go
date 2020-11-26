package ziface

import "net"

type IConnection interface {
	// 启动连接，当前连接开始工作
	Start()
	// 停止连接，结束当前连接状态M
	Stop()
	// 从当前连接获取原始的socket TCPConn
	GetTCPConnection() *net.TCPConn
	// 获取当前连接id
	GetConnID() uint32
	// 获取远程客户端地址信息
	RemoteAddr() net.Addr
	// 直接将Message数据发送数据给远程的TCP客户端
	SendMsg(MsgId uint32, data []byte) error
}

// 定义一个统一处理链接业务的接口
//type HandFunc func(conn *net.TCPConn,bytes []byte,event int) error
