package main

import (
	"io"
	"mxs/log"
	"mxs/gamex/api/tcp/mnet"
	"net"
	"time"
)

func main() {
	log.Debug("Client Test...Start")
	time.Sleep(time.Second*3)

	//conn ,err := mnet.Dial("tcp", "123.56.63.227:7777")
	conn ,err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		log.Error("client dial err", err)
		return
	}
	for {
		// 创建一个封包对象 dp
		dp := mnet.NewDataPack()
		msg, err := dp.Pack(mnet.NewMsgPackage(0, []byte("aaa   0.6 Client Test Message")))
		if err != nil {
			log.Error(" pack msg error")
			return
		}
		_, err = conn.Write(msg)
		if err != nil {
			log.Error("write error err %v", err)
			return
		}
		// 先读出流中的head部分
		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData)
		if err != nil {
			log.Error("read head error")
			break
		}
		// 将headData字节流拆包到msg中
		msgHead, err := dp.UnPack(headData)
		if err != nil {
			log.Error("Unpack err %v", err)
			return
		}
		if msgHead.GetDataLen() > 0 {
			// msg 是有data数据的，需要再次读取data数据
			msg := msgHead.(*mnet.Message)
			msg.Data = make([]byte, msgHead.GetDataLen())
			// 根据datalen 从io中读取字节流
			_, err = io.ReadFull(conn, msg.Data)
			if err != nil{
				log.Error("server unpack data err %v", err)
				return
			}
			log.Debug("===> Recv Msg: ID=%d, len=%d, data=%s", msg.GetMsgId(), msg.GetDataLen(), string(msg.GetData()))
		}
		time.Sleep(1*time.Second)
	}

}