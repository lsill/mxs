package iface

/*
	IRequest 接口：
	实际上是把客户端请求的链接信息 和 请求的数据 包装到了Request里
 */

type IRequest interface {
	GetConnection() IKConnection // 获取请求连接信息
	GetData()	[]byte         // 获取请求消息的数据
	GetMsgTyp() uint32           // 获取请求消息的类型
}