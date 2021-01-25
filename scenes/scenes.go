package main

import (
	"mxs/scenes/core/proc"
	"mxs/util/api/kcp/mnet"
)

func main() {
	s := mnet.NewServer()
	proc.LoadProto(s)
	s.Server()
}
