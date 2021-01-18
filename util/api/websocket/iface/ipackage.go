package iface

import "google.golang.org/protobuf/reflect/protoreflect"

type IPackage interface {
	Reset()
	GetTyp() int32
	GetData()  []byte
	String() string
	ProtoReflect() protoreflect.Message
}