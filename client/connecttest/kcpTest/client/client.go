package main

import (
	"bytes"
	"encoding/binary"
	"github.com/xtaci/kcp-go"
	"mxs/log"
	"mxs/scenes/proto/flat/flatbuffers"
	"mxs/scenes/proto/flat/sample/flatutil"
	"mxs/scenes/proto/flat/sample/strupro"
	"mxs/util/api/kcp/mnet"
	"time"
)



func main() {
	conn , err := kcp.DialWithOptions("localhost:7777", nil, 10, 3)
	if err != nil {
		panic(err)
	}
	builder := flatbuffers.NewBuilder(2000)
	h := flatutil.NewFlatBufferHelper(builder, 32)
	id := h.Pre(builder.CreateString("hello kcp"))
	strupro.TestMessageStart(builder)
	strupro.TestMessageAddTeststr(builder, h.Get(id))
	strupro.TestMessageEnd(builder)
	dp := mnet.NewDataPack()
	bytes :=  builder.Bytes[builder.Head():]
	datalen := len(bytes)
	str := string(bytes)
	log.Debug("str is %v", str)
	//bytes := []byte("test")
	objM := mnet.NewMsgPackage(uint32(0),bytes, int32(datalen))
	msg, err := dp.Pack(objM)
	if err != nil {
		panic(err)
	}

	dq := mnet.NewDataPack()
	obj , err := dq.UnPack(msg)
	if err != nil {
		panic(err)
	}

	log.Debug("%d", obj.GetTyp())
	str = string(obj.GetData())
	log.Debug("%s", str)
	go func() {
		for {
			select {
			case <- time.After(2*time.Second):
				conn.Write(msg)
			}
		}
	}()
	for{
		var buffer = make([]byte, 1024, 1024)
		n, _ := conn.Read(buffer)
		log.Debug("%v", string(buffer[:n]))
		break
	}
}


func main1() {
	dataBuff := bytes.NewBuffer([]byte{})
	var a int32
	a = 5
	err := binary.Write(dataBuff, binary.LittleEndian, a)
	if err != nil {
		panic(err)
	}

	var b int32
	err = binary.Read(dataBuff, binary.LittleEndian, &b)
	if err != nil {
		panic(err)
	}
}



/*EnterPVPGameRspStart(builder)
EnterPVPGameRspAddGameScene(builder, h.Get(GameScene))
EnterPVPGameRspAddAccDatas1(builder, accountInfos1)
EnterPVPGameRspAddAccDatas2(builder, accountInfos2)
EnterPVPGameRspAddBuffData(builder, buffDatas)
EnterPVPGameRspAddPos(builder, int32(p.Place))*/
