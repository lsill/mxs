package mnet

import (
	"mxs/scenes/proto/flat/flatbuffers"
	"mxs/util/api/kcp/iface"
)

type Message struct {
	typ uint32 // 包类型
	builder *flatbuffers.Builder
}

func (m *Message) Typ() uint32{
	return m.typ
}

func (m *Message) Builder() *flatbuffers.Builder {
	return m.builder
}

func NewMsgPackage(typ uint32, data *flatbuffers.Builder) iface.IMessage{
	return &Message{
		typ:     typ,
		builder: data,
	}
}