package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/ziface"
)

/*
	存储一切关Zinx框架的全局参数，供其他模块使用一些参数也可以通过
	用户根据 zinx.json来配置
 */

type GloBalObj struct {
	TcpServer ziface.IServer	// 当前Zinx的全局Server对象
	Host string					// 当前服务器主机IP
	TcpPort int					// 当前服务器主机监听端口号
	Name string					// 当前服务器名称

	MaxPacketSize uint32			// 都需数据包的最大值
	MaxConn		int				// 当前服务器主机允许的最大链接个数
	Version		string 			// 服务器当前版本
}

/*
	定义一个全局的对象
 */
var GlobalObject *GloBalObj

// 读取用户配置文件
func (g *GloBalObj) Reload() {
	data ,err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	// 将json数据解析到struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

func init() {
	// 初始化GlobalObject变量，设置一些默认值
	GlobalObject = &GloBalObj{
		TcpServer:   nil,
		Host:        "",
		TcpPort:     7777,
		Name:        "ZinxServerApp",
		MaxPacketSize: 0,
		MaxConn:     12000,
		Version: "zinx 0.4",
	}
}