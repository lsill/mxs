package iface

// 定义服务器接口
type IServer interface {
	// 启动服务器方法
	Start()
	// 停止服务器方法
	Stop()
	// 开启业务服务器方法
	Server()
	// 路由功能：给当前服务器注册一个路由方法，供客户端链接处理使用
	AddRouter(msgId uint32,router IRouter)
	// 得到连接管理器
	GetConnMgr() IConnManger
	// 设置该Server的连接创建时的Hook函数
	SetOnConnStart(func (connection IConnection))
	// 设置该Server的连接断开时的Hook函数
	SetOnConnStop(func (connection IConnection))
	// 调用连接OnConnStart Hook函数（钩子）
	CallOnConnStart(conn IConnection)
	// 调用连接OnConnStop Hook函数
	CallOnConnStop(conn IConnection)
}

