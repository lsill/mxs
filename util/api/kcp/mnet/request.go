package mnet

import (
	"mxs/util/api/kcp/iface"
)

type Request struct {
	conn iface.IKConnection // 已经和客户端建立好的 链接
	msg  iface.IMessage    // 客户端请求的数据
}

// 获取请求链接信息
func (r *Request) GetConnection() iface.IKConnection {
	return r.conn
}

// 获取请求消息的数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

// 获取请求消息的id
func (r *Request) GetMsgTyp() uint32 {
	return r.msg.GetTyp()
}
