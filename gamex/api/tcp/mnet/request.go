package mnet

import "mxs/gamex/api/tcp/iface"

type Request struct {
	conn iface.IConnection // 已经和客户端建立好的 链接
	msg  iface.IMessage    // 客户端请求的数据
}

// 获取请求链接信息
func (r *Request) GetConnection() iface.IConnection {
	return r.conn
}

// 获取请求消息的数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

// 获取请求消息的id
func (r *Request) GetMsgID() uint32{
	return r.msg.GetMsgId()
}
