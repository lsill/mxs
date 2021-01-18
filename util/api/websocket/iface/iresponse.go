package iface

type IResponse interface {
	SetData([]byte)
	SetPkTyp(int32)
}
