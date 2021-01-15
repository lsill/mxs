package iface

type IMessage interface {
	Type() int
	RealType() int

}
