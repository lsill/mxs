package main

import (
	"fmt"
	"net"
	"time"
	"mxs/log"
)

func main() {
	log.Debug("Clinet Test ...start")
	time.Sleep(3 * time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}
	for {
		_, err := conn.Write([]byte(" v3.0"))
		if err != nil {
			log.Error("write error err", err)
			return
		}
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			log.Error("read buf error")
			return
		}

		log.Debug("Server call back :%s, cnt = %d", buf, cnt)
		time.Sleep(1*time.Second)
	}

}
