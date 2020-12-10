package iface

/*
	封包数据和拆包数据
	直接面向tcp连接中的数据流，为传输数据添加头部信息，用于处理tcp粘包问题
 */

type IDataPack interface {
	GetHeadLen()	uint32				// 获取包头长度的方法
	Pack(msg IMessage) ([]byte, error)	// 封包方法
	UnPack([]byte)(IMessage, error)		// 拆包方法
}


