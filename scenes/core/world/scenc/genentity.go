package scenc

import (
	"mxs/scenes/proto/flat/flatbuffers"
	"mxs/scenes/proto/flat/sample/flatutil"
	"mxs/scenes/proto/flat/sample/strupro"
)

func GenEntityProto(builder *flatbuffers.Builder, entity *Entity) flatbuffers.UOffsetT{
	h := flatutil.NewFlatBufferHelper(builder, 32)
	strupro.EntityStart(builder)
	id := h.Pre(strupro.CreatePosition(builder, entity.X, entity.Y, entity.Z))
	strupro.EntityAddPos(builder, h.Get(id))
	strupro.EntityAddEid(builder, entity.Eid)
	strupro.EntityAddAngle(builder,entity.V)
	return strupro.EntityEnd(builder)
}

func GenPlayersProto(builder *flatbuffers.Builder, entitys []*Player) flatbuffers.UOffsetT {
	_h := flatutil.NewFlatBufferHelper(builder, 32)
	offsets := make([]flatbuffers.UOffsetT, 0,len(entitys))
	for _, play := range entitys {
		offsets = append(offsets, GenEntityProto(builder, play.Entity))
	}
	offset:= _h.CreateUOffsetTArray(strupro.GirdsStartEntityVector, offsets)
	strupro.GirdsStart(builder)
	strupro.GirdsAddEntity(builder, offset)
	return strupro.GirdsEnd(builder)
}

func GenEntitysProto(builder *flatbuffers.Builder, entitys []*Entity) flatbuffers.UOffsetT {
	_h := flatutil.NewFlatBufferHelper(builder, 32)
	strupro.GirdsStart(builder)
	offsets := make([]flatbuffers.UOffsetT, len(entitys))
	for _, play := range entitys {
		offsets = append(offsets, GenEntityProto(builder, play))
	}
	offset:= _h.CreateUOffsetTArray(strupro.GirdsStartEntityVector, offsets)
	strupro.GirdsAddEntity(builder, offset)
	return strupro.GirdsEnd(builder)
}