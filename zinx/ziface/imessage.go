package ziface


/*
	将请求的一个消息封装到Message中，定义一个抽像接口
 */
type IMessage interface {
	GetDataLen()	uint32	// 获取消息数据段长度
	GetMsgId()		uint32	// 获取消息id
	GetData()		[]byte	// 获取消息内容

	SetDataLen(uint32) 	// 设置消息数据段长度
	SetMsgId(uint32)	// 设置消息id
	SetData([]byte)		// 设置消息内容
}


