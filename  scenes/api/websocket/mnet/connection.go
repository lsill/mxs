package mnet

import (
	"errors"
	"github.com/gorilla/websocket"
	logs "mxs/log"
	"mxs/gamex/api/websocket/iface"
	"net"
)

type Connection struct {
	webserver iface.IServer
	conn      *websocket.Conn
	connid    uint32
	isClosed  bool
	// 告知连接已经退出
	ExitBuffChan chan bool
	// 无缓冲通道，俩个goroutine之间通信
	msgChan chan []byte
}


func(c *Connection) Start() {

}

func(c *Connection) Stop() {
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	// TODO	 此处留给回调函数

	c.conn.Close()
	c.ExitBuffChan <- true
	// TODO 此处留给连接管理器删除链接
	close(c.ExitBuffChan)
	close(c.msgChan)
}
func(c *Connection) Conn() *websocket.Conn{
	return c.conn
}

func(c *Connection) ConnId() uint32{
	return c.connid
}

func(c *Connection) RemoteAddr() net.Addr{
	return c.conn.RemoteAddr()
}

func(c *Connection) SendMsg(msgid uint32,data []byte) error{
	if c.isClosed == true {
		return errors.New("connection has closed")
	}
	c.msgChan <- data
	return nil
}

// 读取websocket连接的信息
func (c *Connection) StartReader() {
	logs.Debug("addr %s read is running", c.RemoteAddr())
	defer logs.Debug("addr %s read is stop",c.RemoteAddr())
	defer c.Stop()
	for {
		mt, message, err := c.Conn().ReadMessage()
		if mt != websocket.BinaryMessage {
			logs.Error("message type %s is error", mt)
			return
		}
		if err != nil {
			logs.Error("read err: %v", err)
			break
		}
	}
}