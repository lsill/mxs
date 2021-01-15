package main

import (
	"flag"
	logs "mxs/log"
	"github.com/gorilla/websocket"
	"net/url"
	"os"
	"os/signal"
	"time"
)

var (
	_addr = flag.String("addr", "localhost:8080", "http service address")
)

func main() {
	flag.Parse()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	u := url.URL{
		Scheme:      "ws",
		Host: *_addr,
		Path:        "/echo",
	}
	c, _ , err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		logs.Error("dial: err %v", err)
		return
	}
	defer c.Close()
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message,err := c.ReadMessage()
			if err != nil {
				logs.Error("read err:%s", err)
				return
			}
			logs.Debug("recv message :%s", message)
		}
	}()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <- done:
			return
		case t:=<-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				logs.Error("write err:",err)
				return
			}
		case <- interrupt:
			logs.Release("interrupt")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				logs.Error("write close err:%s", err)
				return
			}
			select {
			case<-done:
			case<-time.After(time.Second):
			}
			return
		}
	}
}