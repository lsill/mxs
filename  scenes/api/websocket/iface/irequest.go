package iface

type IRequest interface {
	GetData() IPacket
	GetSession() ISession
	GetPkTyp() int32
}
