package iface

type IPackage interface {
	Reset()
	GetTyp() int32
	GetData()  []byte
}