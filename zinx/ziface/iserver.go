package ziface

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
}

