package znet

import "C"
import (
	"errors"
	"io"
	"net"
	"zinx/log"
	"zinx/ziface"
)

type Connection struct {
	// 当前连接的socket TCP套接字
	Conn *net.TCPConn
	// 当前连接的ID，可以称作为SessionID,ID全局唯一
	ConnID uint32
	// 当前连接的关闭状态
	isClosed bool
	// 告知该连接已经退出/停止的channel
	ExitBuffChan chan bool
	// 消息管理MsgId和对应处理方法的消息管理模块
	MsgHandler ziface.IMsgHandle
}

// 创建连接的方法
func NewConnecion(conn *net.TCPConn, connId uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		Conn:         conn,
		ConnID:       connId,
		isClosed:     false,
		ExitBuffChan: make(chan bool, 1),
		MsgHandler: msgHandler,
	}
	return  c
}

// 处理conn读数据的Goroutine
func (c *Connection) StartReader() {
	log.Debug("Reader Goroutine is running")
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
			continue
		}

		// 拆包 得到msgid 和 datalen 放在msg中
		msg , err := dp.UnPack(headData)
		if err != nil {
			log.Error("unpack error")
			c.ExitBuffChan<-true
			continue
		}

		// 根据datalen 读取data，放在msg.Data中
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				log.Error("read msg data error:%v", err)
				c.ExitBuffChan<-true
				continue
			}
		}
		msg.SetData(data)

		// 得到当前客户端请求的Request数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		// 从路由Routers 中找到注册绑定Conn对应的Handle
		go func (request ziface.IRequest) {
			// 执行注册的路由方法
			c.MsgHandler.DoMsgHandler(request)
		}(&req)
	}
}

func (c *Connection) Start() {
	go c.StartReader()
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

	// 关闭socket连接
	c.Conn.Close()

	// 通知从缓冲队列读取数据的业务，该连接已经关闭
	c.ExitBuffChan <- true

	// 关闭该连接的全部管道
	close(c.ExitBuffChan)
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
	if _, err := c.Conn.Write(msg); err != nil {
		log.Error("Write msg id:%d error", msgId)
		c.ExitBuffChan<- true
		return errors.New("conn Write error")
	}

	return nil
}