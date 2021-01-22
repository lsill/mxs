package iface

import "mxs/scenes/proto/flat/flatbuffers"

/*
	将请求的一个消息封装到Message中，定义一个抽像接口
 */
type IMessage interface {
	Typ()	uint32 	// 获取消息类型
	Builder()	*flatbuffers.Builder	// 获取消息内容
}


