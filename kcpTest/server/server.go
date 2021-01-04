package main

import (
	"github.com/xtaci/kcp-go"
	"io"
	"mxs/log"
	"net"
)

func main() {
	lis,err := kcp.ListenWithOptions(":9999", nil, 10, 3)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := lis.AcceptKCP()
		if err != nil {
			panic(err)
		}
		go func (conn net.Conn) {
			var buffer = make([]byte, 1024, 1024)
			for {
				n, err:= conn.Read(buffer)
				if err != nil {
					if err == io.EOF {
						break
					}
					log.Debug("%v", err)
					break
				}
				log.Debug("receive from client: %v", string(buffer[:n]))
				conn.Write([]byte("respone kcp!"))
			}
		}(conn)
	}
}
