package mnet

import "C"
import (
	"errors"
	"github.com/xtaci/kcp-go"
	"io"
	"mxs/gamex/api/tcp/iface"
	"mxs/gamex/utils"
	"mxs/log"
	"net"
	"sync"
)

type Connection struct {
	// tcp 服务
	TcpServer iface.IServer
	// 当前连接的socket TCP套接字
	Conn *net.TCPConn
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
func NewConnecion(server iface.IServer,conn *net.TCPConn, connId uint32, msgHandler iface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer: server,	// 将隶属的server传递进来
		Conn:         conn,
		ConnID:       connId,
		isClosed:     false,
		ExitBuffChan: make(chan bool, 1),
		MsgHandler: msgHandler,
		msgChan: make(chan []byte),
		msgBuffChan: make(chan []byte, utils.GloUtil.MaxMsgChanLen),
		property: make(map[string]interface{}), // 初始化连接属性map
	}
	// 将新创建的Conn添加到连接管理中
	c.TcpServer.GetConnMgr().Add(c) // 将当前心新创建的连接添加到ConnManager中
	return  c
}

// 处理conn读数据的Goroutine
func (c *Connection) StartReader() {
	log.Debug("Reader %v is running", c.RemoteAddr())
	defer log.Debug("%s conn reader exit!",c.RemoteAddr().String())
	defer c.Stop()
	for {
		// 创建拆包解包的对象
		dp := NewDataPack()

		// 读取客户端的msg head
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil{
			log.Error("read msg head error!")
			c.ExitBuffChan<-true
			break
		}

		// 拆包 得到msgid 和 datalen 放在msg中
		msg , err := dp.UnPack(headData)
		if err != nil {
			log.Error("unpack error")
			c.ExitBuffChan<-true
			break
		}

		// 根据datalen 读取data，放在msg.Data中
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				log.Error("read msg data error:%v", err)
				continue
			}
		}
		msg.SetData(data)

		// 得到当前客户端请求的Request数据
		req := Request{
			conn: c,
			msg:  msg,
		}
		if utils.GloUtil.MaxWorkerTaskLen > 0{
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
func (c *Connection) StartWriter() {
	log.Debug("Writer %v is running", c.RemoteAddr())
	defer log.Debug("conn Writer %v exit!", c.RemoteAddr())
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


func (c *Connection) Start() {
	// 1. 开启从客户端读取数据流程的goroutine
	go c.StartReader()
	// 2.开启写回客户端数据流程的goroutine
	go c.StartWriter()

	// 按照传进来创建连接时需要处理的业务，执行钩子方法
	c.TcpServer.CallOnConnStart(c)

	for {
		select {
		case <- c.ExitBuffChan:
		// 得到退出消息,不在阻塞
			return
		}
	}
}

// 获取远程客户端的地址信息
func (c *Connection) RemoteAddr() net.Addr{
	return c.Conn.RemoteAddr()
}

func (c *Connection) Stop() {
	// 1.如果当前连接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	// TODO Connection Stop() 如果用户注册了该连接的关闭回调业务，那么此刻应该显示调用
	log.Debug("conn stop callonconnstop")
	c.TcpServer.CallOnConnStop(c)


	// 关闭socket连接
	c.Conn.Close()

	// 通知从缓冲队列读取数据的业务，该连接已经关闭
	c.ExitBuffChan <- true
	// 将连接从连接管理器中删除
	c.TcpServer.GetConnMgr().Remove(c)
	// 关闭该连接的全部管道
	close(c.ExitBuffChan)
	close(c.msgBuffChan)
	close(c.msgChan)
}

// 获取当前连接id
func (c *Connection) GetConnID() uint32{
	return c.ConnID
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 直接将Message数据发送数据给远程的TCP客户端
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}

	// 将data封包，并且发送
	dp := NewDataPack()
	msg , err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		log.Error("Pack error msg id = %d", msgId)
		return errors.New("Pack error msg")
	}

	// 写回客户端
	c.msgChan <- msg	// 将之前直接写会给conn.writer的方法 改为发送给Channel 供writer读取

	return nil
}


func (c *Connection) SendBuffMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}

	// 将data封包，并且发送
	dp := NewDataPack()
	msg ,err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		log.Error("buff Pack error msg id = %d", msgId)
		return errors.New("buff Pack error msg")
	}
	c.msgBuffChan <- msg 	// 写回客户端
	return nil
}

// 设置连接属性
func (c *Connection)SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}

// 获取链接属性
func (c *Connection)GetProperty(key string) (interface{}, error) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New(key + " property not found")
	}
}

// 删除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
}

type KConnection struct {
	// kcp服务器
	KcpServer iface.IServer
	// 当前连接的Kcp会话
	Session *kcp.UDPSession
	// 当前连接的会话id，全局唯一id
	SessionId uint32
	// 回写通道
	msgChan chan []byte
}

func (c *KConnection) Start() {

}

func (c *KConnection) Stop() {

}


func (c *KConnection) GetUdpSession() *kcp.UDPSession {
	return  c.Session
}

func (c *KConnection) GetKcpId() uint32 {
	return c.SessionId
}

func (c *KConnection) GetConnectionAddr() net.Addr {
	return c.Session.RemoteAddr()
}

func (c *KConnection) SendMsg(msgId uint32, data []byte) (int, error) {
	// 将data封包，并且发送
	dp := NewDataPack()
	msg , err := dp.Pack(NewMsgPackage(msgId, data))
	n, err := c.Session.Write(data)
	if err != nil {
		return n, err
	}
	c.msgChan <- msg
	return n, nil
}




