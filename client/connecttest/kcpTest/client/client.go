package main

import (
	"bytes"
	"encoding/binary"
	"github.com/xtaci/kcp-go"
	"mxs/log"
	"mxs/scenes/proto/flat/sample/strupro"
	"mxs/util/api/kcp/mnet"
	"time"
)



func main() {
	conn , err := kcp.DialWithOptions("localhost:7777", nil, 10, 3)
	if err != nil {
		panic(err)
	}
	//builder := flatbuffers.NewBuilder(2000)
	//h := flatutil.NewFlatBufferHelper(builder, 32)
	//id := h.Pre(builder.CreateString("hello kcp"))
	//strupro.TestMessageStart(builder)
	//strupro.TestMessageAddTeststr(builder, h.Get(id))
	//strupro.TestMessageEnd(builder)
	dp := mnet.NewDataPack()
	//bytes :=  builder.Bytes[builder.Head():]
	bytes := make([]byte, 0)
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

	str = string(obj.GetData())
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
		conn.Read(buffer) //_, _ :=
		db := mnet.NewDataPack()
		getmsg,err := db.UnPack(buffer)
		if err != nil {
			log.Error("unpack error")
		}
		if getmsg.GetTyp() == 1 {
			gird := strupro.GetRootAsGirds(getmsg.GetData(), 0)
			log.Debug("get len is %d", gird.EntityLength())
			if gird.EntityLength() > 0{
				entity := &strupro.Entity{}
				log.Debug("get eid is %v", gird.Entity(entity,0))
				log.Debug("eid is %d",entity.Eid())
			}

		}
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
