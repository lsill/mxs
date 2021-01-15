package mnet

type Packet struct {
	typ int32
	message []byte
}

func (rcv *Packet) Typ() int32 {
	return rcv.typ
}

func (rcv *Packet) Data() []byte {
	return rcv.message
}

func NewPacket(typ int32,message []byte) *Packet {
	return &Packet{
		message: message,
		typ: typ,
	}
}
