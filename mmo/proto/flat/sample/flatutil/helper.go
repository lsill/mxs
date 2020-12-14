package flatutil

import "mxs/mmo/proto/flat/flatbuffers"

func GetNewBuilder() *flatbuffers.Builder {
	return flatbuffers.NewBuilder(20480)
}

type flatbufferPreBuildHelper struct {
	offsets []flatbuffers.UOffsetT
	builder *flatbuffers.Builder
}

func NewFlatBufferHelper(b *flatbuffers.Builder, cap int) *flatbufferPreBuildHelper {
	return &flatbufferPreBuildHelper{
		builder: b,
		offsets: make([]flatbuffers.UOffsetT, 0, cap),
	}
}

func (f *flatbufferPreBuildHelper) Pre(o flatbuffers.UOffsetT) int {
	f.offsets = append(f.offsets, o)
	return len(f.offsets) - 1
}

func (f *flatbufferPreBuildHelper) PreStringAr(ar []string) []int {
	res := make([]int, 0, len(ar))
	for i := 0; i < len(ar); i++ {
		f.offsets = append(f.offsets, f.builder.CreateString(ar[i]))
		res = append(res, len(f.offsets)-1)
	}
	return res
}

func (f *flatbufferPreBuildHelper) PreStringArForUOffsetT(ar []string) []flatbuffers.UOffsetT {
	res := make([]flatbuffers.UOffsetT, 0, len(ar))
	var o flatbuffers.UOffsetT
	for i := 0; i < len(ar); i++ {
		o = f.builder.CreateString(ar[i])
		f.offsets = append(f.offsets, o)
		res = append(res, o)
	}
	return res
}

func (f *flatbufferPreBuildHelper) Get(idx int) flatbuffers.UOffsetT {
	return f.offsets[idx] // 如果对不上直接崩溃
}

func (f *flatbufferPreBuildHelper) GetAr(idxs []int) []flatbuffers.UOffsetT {
	res := make([]flatbuffers.UOffsetT, 0, len(idxs))
	for i := 0; i < len(idxs); i++ {
		res = append(res, f.Get(idxs[i]))
	}
	return res[:]
}

type flatbufVecStartFunc func(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT

func (f *flatbufferPreBuildHelper) CreateUInt32Array(fc flatbufVecStartFunc, ar []uint32) flatbuffers.UOffsetT {
	numElems := len(ar)
	fc(f.builder, numElems)
	for i := numElems - 1; i >= 0; i-- {
		f.builder.PrependUint32(ar[i])
	}
	return f.builder.EndVector(numElems)
}

func (f *flatbufferPreBuildHelper) CreateIntArray(fc flatbufVecStartFunc, ar []int) flatbuffers.UOffsetT {
	numElems := len(ar)
	fc(f.builder, numElems)
	for i := numElems - 1; i >= 0; i-- {
		f.builder.PrependInt32(int32(ar[i]))
	}
	return f.builder.EndVector(numElems)
}

func (f *flatbufferPreBuildHelper) CreateStringArray(fc flatbufVecStartFunc, ar []string) flatbuffers.UOffsetT {
	return f.CreateUOffsetTArray(fc, f.PreStringArForUOffsetT(ar))
}

func (f *flatbufferPreBuildHelper) CreateUOffsetTArray(fc flatbufVecStartFunc, ar []flatbuffers.UOffsetT) flatbuffers.UOffsetT {
	numElems := len(ar)
	fc(f.builder, numElems)
	for i := numElems - 1; i >= 0; i-- {
		f.builder.PrependUOffsetT(ar[i])
	}
	return f.builder.EndVector(numElems)
}