package main

import (
	"github.com/xtaci/kcp-go"
	"mxs/log"
)

func main() {
	conn , err := kcp.DialWithOptions("localhost:9999", nil, 10, 3)
	if err != nil {
		panic(err)
	}
	conn.Write([]byte("hello kcp!"))
	for{
		var buffer = make([]byte, 1024, 1024)
		n, _ := conn.Read(buffer)
		log.Debug("%v", string(buffer[:n]))
		break
	}
}
