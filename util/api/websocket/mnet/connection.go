package mnet

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	logs "mxs/log"
	"mxs/util/api/websocket/iface"
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

func NewConnecton(server iface.IServer, conn *websocket.Conn,connid uint32, msghandler iface.IMsgHandle) *Connection {
	return &Connection{
		webserver:    server,
		conn:         conn,
		connid:       connid,
		isClosed:     false,
		ExitBuffChan: make(chan bool, 1),
		msgChan:      make(chan []byte),
	}
}

func(c *Connection) Start() {
	go c.StartReader()
	go c.StartWriter()
	// TODO 此处留着处理hook

	for {
		select {
		case <- c.ExitBuffChan:
			return
		}
	}
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
		data := &Package{}
		err = proto.Unmarshal(message, data)
		if err != nil {
			logs.Error("reader: unmarshal err")
			return
		}
		logs.Debug("reader data:%v", data)
		req := Request{
			conn: c,
			pk:   data,
		}
		logs.Release("req is %v", req)
	}
}

func (c *Connection) StartWriter() {
	logs.Debug("Writer %v is Running", c.RemoteAddr())
	defer logs.Debug("conn writer %v exit!", c.RemoteAddr())
	for {
		select {
		case data, ok := <-c.msgChan:
			if ok {
				err := c.Conn().WriteMessage(websocket.BinaryMessage, data)
				if err != nil{
					logs.Error("Send Data error, ConnWriter exit!")
					return
				}
			} else {
				logs.Error("msgChan is cloesd")
				break
			}
		case <- c.ExitBuffChan:
			return
		}
	}
}