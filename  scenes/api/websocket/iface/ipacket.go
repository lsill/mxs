package iface

type IPacket interface {
	Typ() int32
	Data() []byte
	Pack(msg )
}
