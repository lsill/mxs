package main

import "mxs/util/api/kcp/mnet"

func main() {
	s := mnet.NewServer()
	s.Server()
}
