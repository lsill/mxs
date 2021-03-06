// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package strupro

import (
	flatbuffers "mxs/scenes/proto/flat/flatbuffers"
)

type Girds struct {
	_tab flatbuffers.Table
}

func GetRootAsGirds(buf []byte, offset flatbuffers.UOffsetT) *Girds {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Girds{}
	x.Init(buf, n+offset)
	return x
}

func GetSizePrefixedRootAsGirds(buf []byte, offset flatbuffers.UOffsetT) *Girds {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &Girds{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func (rcv *Girds) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Girds) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Girds) Entity(obj *Entity, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *Girds) EntityLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func GirdsStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func GirdsAddEntity(builder *flatbuffers.Builder, entity flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(entity), 0)
}
func GirdsStartEntityVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func GirdsEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
