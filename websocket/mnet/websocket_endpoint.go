package mnet

import (
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
)

type WebSocketEndPoint struct {
	ws *websocket.Conn
	output chan []byte
	mtype int
	ConnId uint32
}

func NewWebSocketEndPoint(ws *websocket.Conn, bin bool) *WebSocketEndPoint{
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

func (we *WebSocketEndPoint) GetConnId() {
	return we.ws.
}

func (we *WebSocketEndPoint) Terminate() {
	logs.Info("websocket")
}

func (we *WebSocketEndPoint) Output() {

}

func (we *WebSocketEndPoint) Send() {

}

func (we *WebSocketEndPoint) StartReading() {

}
