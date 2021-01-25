package iface


/*
	将请求的一个消息封装到Message中，定义一个抽像接口
 */
type IMessage interface {
	GetDataLen() uint32
	GetMsgId() uint32
	GetData() []byte
	SetData([]byte)
}
