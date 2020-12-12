package mnet

type Message struct {
	Id uint32	// 消息id
	DataLen	 uint32	// 消息的长度
	Data	[]byte	// 消息内容
}

// 创建一个消息包
func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		Id:      id,
		DataLen: uint32(len(data)),
		Data:   data,
	}
}

// 获取消息数据段长度
func (msg *Message) GetDataLen() uint32{
	return msg.DataLen
}

// 获取消息id
func (msg *Message) GetMsgId() uint32 {
	return msg.Id
}

// 获取消息内容
func (msg *Message) GetData() []byte {
	return msg.Data
}

// 设置消息段长度
func (msg *Message) SetDataLen(len uint32) {
	msg.DataLen = len
}

// 设置消息id
func (msg *Message) SetMsgId(id uint32) {
	msg.Id = id
}

// 设置消息内容
func (msg *Message) SetData(data []byte) {
	msg.Data = data
}
