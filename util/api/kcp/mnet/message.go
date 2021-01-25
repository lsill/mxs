package mnet

type Message struct {
	Typ uint32 // 包类型
	DataLen int32	// 数据长度
	Data []byte
}

func (m *Message) GetTyp() uint32{
	return m.Typ
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) GetDataLen() int32 {
	return m.DataLen
}

func NewMsgPackage(typ uint32, data []byte, datalen int32) *Message{
	return &Message{
		Typ:     typ,
		Data: data,
		DataLen: datalen,
	}
}