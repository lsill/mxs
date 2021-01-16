package mnet

import (
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"io"
	"io/ioutil"
)

type WebSocketEndPoint struct {
	ws *websocket.Conn
	output chan []byte
	mtype int
}

func NewWebSocketEndPoint(ws *websocket.Conn, bin bool) *WebSocketEndPoint {
	endpoint := &WebSocketEndPoint{
		ws:     ws,
		output: make(chan []byte),
		mtype:  websocket.TextMessage,
	}
	if bin {
		endpoint.mtype = websocket.BinaryMessage
	}
	return endpoint
}


func (we *WebSocketEndPoint) Terminate() {
	logs.Info("terminate websocket connection")
}

func (we *WebSocketEndPoint) Output() chan []byte{
	return we.output
}

func (we *WebSocketEndPoint) Send(msg []byte) bool{
	w, err := we.ws.NextWriter(we.mtype)
	if err == nil {
		_, err = w.Write(msg)
	}
	w.Close()
	if err != nil {
		logs.Error("websocket can't send:%v", err)
		return false
	}
	return true
}

func (we *WebSocketEndPoint) StartReading() {
	go we.read_frames()
}

func (we *WebSocketEndPoint) read_frames() {
	for {
		mtype,rd,err := we.ws.NextReader()
		if err != nil {
			logs.Error("websocket can't receive:%v", err)
			break
		}
		if mtype != we.mtype{
			logs.Warn("websocket received message of type is not except, ignoring...")
		}
		p, err := ioutil.ReadAll(rd)
		if err != nil && err != io.EOF {
			logs.Warn("websocket can't received message:%v", err)
			break
		}
		switch mtype {
		case websocket.TextMessage:
			we.output<-append(p, '\n')
		case websocket.BinaryMessage:
			we.output<-p
		default:
			logs.Error("websocket received message of unknown type:%d", mtype)
		}
	}
	close(we.output)
}
