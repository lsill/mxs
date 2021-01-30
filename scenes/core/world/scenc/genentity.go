package scenc

import (
	"mxs/scenes/proto/flat/flatbuffers"
	"mxs/scenes/proto/flat/sample/flatutil"
	"mxs/scenes/proto/flat/sample/strupro"
)

func GenEntityProto(builder *flatbuffers.Builder, entity *Entity) flatbuffers.UOffsetT{
	h := flatutil.NewFlatBufferHelper(builder, 32)
	id := h.Pre(strupro.CreatePosition(builder, entity.X, entity.Y, entity.Z))
	strupro.EntityStart(builder)
	strupro.EntityAddPos(builder, h.Get(id))
	strupro.EntityAddEid(builder, entity.Eid)
	strupro.EntityAddAngle(builder,entity.V)
	return strupro.EntityEnd(builder)
}

func GenPlayersProto(builder *flatbuffers.Builder, entitys []*Player) flatbuffers.UOffsetT {
	_h := flatutil.NewFlatBufferHelper(builder, 32)
	offset := make([]flatbuffers.UOffsetT, 0,len(entitys))
	for _, play := range entitys {
		offset = append(offset, GenEntityProto(builder, play.Entity))
	}
	offsets:= _h.CreateUOffsetTArray(strupro.GirdsStartEntityVector, offset)
	strupro.GirdsStart(builder)
	strupro.GirdsAddEntity(builder, offsets)
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