package iface

type IRequest interface {
	GetData() IPackage
	GetConn() IConnection
	GetPkTyp() int32
}
