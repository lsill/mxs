package main

import (
	"io"
	"mxs/log"
	"mxs/util/api/tcp/mnet"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		return
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error("server accept err:%v", err)
		}
		go func (conn net.Conn) {
			dp := mnet.NewDataPack()
			for {
				// 1 先读出流中的head部分
				headData := make([]byte, dp.GetHeadLen())
				_, err := io.ReadFull(conn, headData) // ReadFull 会把msg填充满为止
				if err != nil {
					log.Error("read head err:%v", err)
					break
				}
				// 将headData字节流 拆包到msg中
				msgHead, err := dp.UnPack(headData)
				if err != nil {
					log.Error("server unpack err")
					return
				}
				if msgHead.GetDataLen() > 0 {
					// msg 是有data数据的，需要再次读取data数据
					msg := msgHead.(*mnet.Message)
					msg.Data = make([]byte,msgHead.GetDataLen())

					// 根据dataLen从io从读取字节流
					_, err := io.ReadFull(conn, msg.Data)
					if err != nil {
						log.Error("server unpack data err:", err)
						return
					}
					log.Debug("==> Recv Msg: ID=%d, len=%d, data=%v", msg.Id, msg.DataLen, string(msg.Data))
				}
			}
		}(conn)
	}
}
