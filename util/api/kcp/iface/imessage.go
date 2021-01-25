package iface

/*
	将请求的一个消息封装到Message中，定义一个抽像接口
 */
type IMessage interface {
	GetTyp()	uint32 	// 获取消息类型
	GetData() []byte// 获取消息内容
	GetDataLen() int32
}


