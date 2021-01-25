package mnet

import (
	"errors"
	"github.com/xtaci/kcp-go"
	"io"
	"mxs/log"
	"mxs/util"
	"mxs/util/api/kcp/iface"
	"net"
	"sync"
)

type KConnection struct {
	// tcp 服务
	Server iface.IServer
	// 当前连接的socket TCP套接字
	Conn *kcp.UDPSession
	// 当前连接的ID，可以称作为SessionID,ID全局唯一
	ConnID uint32
	// 当前连接的关闭状态
	isClosed bool
	// 告知该连接已经退出/停止的channel
	ExitBuffChan chan bool
	// 消息管理MsgId和对应处理方法的消息管理模块
	MsgHandler iface.IMsgHandle
	// 无缓冲通道，用于读、写两个goroutine之间的消息通信
	msgChan	chan []byte
	// 有缓冲通道，用于读、写两个goroutine之间的消息通信
	msgBuffChan chan[]byte

	// 连接属性
	property map[string]interface{}
	//保护连接属性修改的锁
	propertyLock sync.RWMutex
}

// 创建连接的方法
func NewConnecion(server iface.IServer,conn *kcp.UDPSession, connId uint32, msgHandler iface.IMsgHandle) *KConnection {
	c := &KConnection{
		Server: server,	// 将隶属的server传递进来
		Conn:         conn,
		ConnID:       connId,
		isClosed:     false,
		ExitBuffChan: make(chan bool, 1),
		MsgHandler: msgHandler,
		msgChan: make(chan []byte),
		msgBuffChan: make(chan []byte, util.GloUtil.MaxMsgChanLen),
		property: make(map[string]interface{}), // 初始化连接属性map
	}
	// 将新创建的Conn添加到连接管理中
	c.Server.GetConnMgr().Add(c) // 将当前心新创建的连接添加到ConnManager中
	return  c
}

// 处理conn读数据的Goroutine
func (c *KConnection) StartReader() {
	log.Debug("Reader %v is running", c.GetConnectionAddr())
	defer log.Debug("%s conn reader exit!",c.GetConnectionAddr().String())
	defer c.Stop()
	for {
		// 创建拆包解包的对象
		dp := NewDataPack()

		// 读取客户端的msg head
		var buffer = make([]byte, 1024, 1024)
		 _, err := c.Conn.Read(buffer)
		 if err != nil {
		 	if err == io.EOF{
				c.ExitBuffChan<-true
		 		break
			}
			log.Debug("kcp read err %v", err)
			 c.ExitBuffChan<-true
		 	break
		 }
		// 拆包 得到msgid 和 datalen 放在msg中
		msg , err := dp.UnPack(buffer)
		log.Debug("get msg is %v", string(msg.GetData()))
		log.Debug("get msg typ is %v", msg.GetTyp())
		log.Debug("get msg len is %v",msg.GetDataLen())
		if err != nil {
			log.Error("unpack error")
			c.ExitBuffChan<-true
			break
		}

		// 得到当前客户端请求的Request数据
		req := Request{
			conn: c,
			msg:  msg,
		}
		if util.GloUtil.MaxWorkerTaskLen > 0{
			// 已经启动工作池机制，将消息交给worker处理
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			// 从路由Routers 中找到注册绑定Conn对应的Handle
			// 执行注册的路由方法
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

/*
	写消息Goroutine,用户将数据发送给客户端
 */
func (c *KConnection) StartWriter() {
	log.Debug("Writer %v is running", c.GetConnectionAddr())
	defer log.Debug("conn Writer %v exit!", c.GetConnectionAddr())
	//defer c.Stop()
	for {
		select {
			case data, ok :=<-c.msgChan:
				// 有数据要写给客户端
				if ok {
					if _, err := c.Conn.Write(data); err != nil {
						log.Error("Send Data error %v, ConnWriter exit!")
						return
					}
				} else {
					log.Error("msgChan is Closed")
					break
				}

			case data, ok := <-c.msgBuffChan:
				if ok {
					if _, err := c.Conn.Write(data); err != nil {
						log.Error("Send Data error %v, ConnWriterBuff eixt!")
						return
					}
				}else {
					log.Error("msgBuffChan is Closed")
					break
				}


			case <- c.ExitBuffChan:
			// conn已经关闭
				return
		}
	}
}


func (c *KConnection) Start() {
	// 1. 开启从客户端读取数据流程的goroutine
	go c.StartReader()
	// 2.开启写回客户端数据流程的goroutine
	go c.StartWriter()

	// 按照传进来创建连接时需要处理的业务，执行钩子方法
	c.Server.CallOnConnStart(c)

	for {
		select {
		case <- c.ExitBuffChan:
		// 得到退出消息,不在阻塞
			return
		}
	}
}

// 获取远程客户端的地址信息
func (c *KConnection) RemoteAddr() net.Addr{
	return c.Conn.RemoteAddr()
}

func (c *KConnection) KConnection() {
	// 1.如果当前连接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	// TODO Connection Stop() 如果用户注册了该连接的关闭回调业务，那么此刻应该显示调用
	log.Debug("conn stop callonconnstop")
	c.Server.CallOnConnStop(c)


	// 关闭socket连接
	c.Conn.Close()

	// 通知从缓冲队列读取数据的业务，该连接已经关闭
	c.ExitBuffChan <- true
	// 将连接从连接管理器中删除
	c.Server.GetConnMgr().Remove(c)
	// 关闭该连接的全部管道
	close(c.ExitBuffChan)
	close(c.msgBuffChan)
	close(c.msgChan)
}

// 获取当前连接id
func (c *KConnection) GetConnID() uint32{
	return c.ConnID
}

func (c *KConnection) GetKCPConnection() *kcp.UDPSession {
	return c.Conn
}

// 直接将Message数据发送数据给远程的TCP客户端
func (c *KConnection) SendMsg(msgId uint32, data []byte, datalen int32) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}

	// 将data封包，并且发送
	dp := NewDataPack()
	msg , err := dp.Pack(NewMsgPackage(msgId, data, datalen))
	if err != nil {
		log.Error("Pack error msg id = %d", msgId)
		return errors.New("Pack error msg")
	}

	// 写回客户端
	c.msgChan <- msg	// 将之前直接写会给conn.writer的方法 改为发送给Channel 供writer读取

	return nil
}


func (c *KConnection) SendBuffMsg(msgId uint32, data []byte,datalen int32) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}

	// 将data封包，并且发送
	dp := NewDataPack()
	msg ,err := dp.Pack(NewMsgPackage(msgId, data, datalen))
	if err != nil {
		log.Error("buff Pack error msg id = %d", msgId)
		return errors.New("buff Pack error msg")
	}
	c.msgBuffChan <- msg 	// 写回客户端
	return nil
}

// 设置连接属性
func (c *KConnection)SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}

// 获取链接属性
func (c *KConnection)GetProperty(key string) (interface{}, error) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New(key + " property not found")
	}
}

// 删除链接属性
func (c *KConnection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
}



func (c *KConnection) GetUdpSession() *kcp.UDPSession {
	return  c.Conn
}

func (c *KConnection) GetKcpId() uint32 {
	return c.ConnID
}

func (c *KConnection) GetConnectionAddr() net.Addr {
	return  c.RemoteAddr()
}


func (c *KConnection) Stop() {

}


