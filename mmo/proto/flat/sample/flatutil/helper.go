package flatutil

import "mxs/mmo/proto/flat/flatbuffers"

func GetNewBuilder() *flatbuffers.Builder {
	return flatbuffers.NewBuilder(20480)
}