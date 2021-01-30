package main

import (
	"mxs/scenes/core/proc"
	"mxs/util/api/kcp/mnet"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	s := mnet.NewServer()
	proc.LoadProto(s)
	go func() {
		http.ListenAndServe(":8080", nil)
	}()
	s.Server()
}
