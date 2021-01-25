package mnet

import (
	"bytes"
	"encoding/binary"
	"mxs/util/api/kcp/iface"
)

// 封包拆包实力，暂时不需要成员
type DataPack struct {

}

type SliceMock struct {
	Addr uintptr
	Len  int
	Cap  int
}

// 封包拆包实例初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

// 获取包头长度方法
func(dp *DataPack) GetHeadLen() int64{
	// Id uint32（4字节） + DataLen uint32(4字节)
	return 8
}

// 封包方法（压缩数据） 此处留下疑问 为什么要用小端序
func (dp *DataPack) Pack(msg iface.IMessage)([]byte, error){
	// 创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetTyp()); err != nil {
		return nil, err
	}

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// 拆包方法(解压数据)
func (dp *DataPack) UnPack(binaryData []byte) (iface.IMessage, error){
	// 创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)
	msg := &Message{}

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Typ); err != nil {
		return nil, err
	}

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	msg.Data = make([]byte, msg.DataLen)
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Data); err != nil{
		return nil, err
	}

	/*if (util.GloUtil.MaxPacketSize > 0 && msg.DataLen > util.GloUtil.MaxPacketSize) {
		return nil, errors.New("Too large msg data recived!")
		log.Error("Too large msg data recived!")
	}*/

	return msg, nil
}