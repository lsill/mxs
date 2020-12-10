package utils

import (
	"encoding/json"
	"io/ioutil"
	"mxs/api/iface"
)

/*
	存储一切关Zinx框架的全局参数，供其他模块使用一些参数也可以通过
	用户根据 api.json来配置
 */

type GloBalObj struct {
	TcpServer iface.IServer // 当前Zinx的全局Server对象
	Host      string        // 当前服务器主机IP
	TcpPort   int           // 当前服务器主机监听端口号
	Name      string        // 当前服务器名称

	MaxPacketSize uint32			// 都需数据包的最大值
	MaxConn		int				// 当前服务器主机允许的最大链接个数
	Version		string 			// 服务器当前版本
	MaxWorkerTaskLen uint32		// 当前工作worker池的数量 如果为0的话就不开启工作池机制

	MaxMsgChanLen uint32		// 最大缓冲通道长度

	ConFilePath string		// 配置文件路径
}

/*
	定义一个全局的对象
 */
var GloUtil *GloBalObj

// 读取用户配置文件
func (g *GloBalObj) Reload() {
	data ,err := ioutil.ReadFile("conf/api.json")
	if err != nil {
		panic(err)
	}
	// 将json数据解析到struct中
	err = json.Unmarshal(data, &GloUtil)
	if err != nil {
		panic(err)
	}
}

func init() {
	// 初始化GloUtil变量，设置一些默认值
	GloUtil = &GloBalObj{
		TcpServer:   nil,
		Host:        "",
		TcpPort:     7777,
		Name:        "ZinxServerApp",
		MaxPacketSize: 0,
		MaxConn:     12000,
		Version: "api 0.4",
		MaxWorkerTaskLen:100,
		MaxMsgChanLen:1000,
	}
	// 从配置文件中加载一些配置参数
	GloUtil.Reload()
}