package iface

type IRequest interface {
	GetData() []byte
	GetConn() IConnection
	GetPkTyp() int32
}
