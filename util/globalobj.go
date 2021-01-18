package util

import (
	"encoding/json"
	"io/ioutil"
	logs "mxs/log"
	"os"
	"strings"
	"time"
)

/*
	存储一切关框架的全局参数，供其他模块使用一些参数也可以通过
	用户根据 api.json来配置
 */

type GloBalObj struct {
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
var WebsocketObj *WebsocketConfig

// 读取用户配置文件
func (g *GloBalObj) Reload() {
	_path, _ := os.Getwd()
	_path = strings.Replace(_path,"gamex", "", 1)
	path := _path+ "/util/api/conf/api.json"
	logs.Release("path is %s", path)
	data ,err := ioutil.ReadFile(path)
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
		Host:        "",
		TcpPort:     7777,
		Name:        "ServerApp",
		MaxPacketSize: 0,
		MaxConn:     12000,
		Version: "api 0.4",
		MaxWorkerTaskLen:100,
		MaxMsgChanLen:1000,
	}
	// 从配置文件中加载一些配置参数
	GloUtil.Reload()
}

/*func (g *WebsocketConfig) Reload() {
	data, err := ioutil.ReadFile("conf/websocket.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &WebsocketObj)
	if err != nil {
		panic(err)
	}
}*/

type WebsocketConfig struct {
	// 服务器开启时间（用于开发者控制台缓存）
	StartupTime time.Time
	// 要执行的命令。
	CommandName string
	// 传递给命令其他参数。
	CommandArgs []string
	//  传递给SERVER_SOFTWARE环境变量的值(例：websocket/1.2.3)
	ServerSoftware string
	// 毫秒开始发送信号
	CloseMs uint
	// 完成握手的时间（默认为1500ms）
	HandshakeTimeOut time.Duration

	// 一下是相关设定

	// 使用二进制通信（以块的形式发送数据，这些数据是从进程中读取的）
	Binary bool
	// 对主机名执行反向DNS查找（有用，但速度较慢）。
	ReverseLookup bool
	// websocketd与--ssl一起使用，这意味着正在使用TLS
	Ssl 	bool
	// websocket脚本的基本目录。
	ScripDir string
	// 是否正在使用脚本运行目录
	UsingScriptDir bool

	// 如果设置，静态文件将通过HTTP从该目录提供。
	StaticDir string
	// 如果设置，则将通过HTTP从该目录提供CGI脚本。
	CgiDir string
	// 启用开发者控制台。这将禁用StaticDir和CgiDir。
	DevConsole bool
	// websocket升级允许的原始地址列表
	AllowOrigins []string
	// 如果设置，则仅要求从同一来源执行websocket升级。
	SameOrigin bool

	Headers        []string
	HeadersWs      []string
	HeadersHTTP    []string

	// 创建环境
	Env []string	// 要传递给进程的其他环境变量（"key=value"）。
	ParentEnv []string	// 在清理子进程之前，将其保留在os.Environ（）中。
}