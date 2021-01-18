package req

import (
	"github.com/golang/protobuf/proto"
	"github.com/prometheus/common/log"
	logs "mxs/log"
	"mxs/util/api/websocket/iface"
	"mxs/util/api/websocket/mnet"
	pb "mxs/gamex/proto/protoc/pb"
)

type HelloRouter struct {
	mnet.BaseRouter
}

func (this *HelloRouter) PreHandler(req iface.IRequest) {
	logs.Debug("Call Router preHandle")
}

func (this *HelloRouter) Handler(req iface.IRequest) {
	log.Debug("Call PingRouter Handle")
	hello := &pb.HelloResp{}
	err := proto.Unmarshal(req.GetData(), hello)
	if err != nil {
		logs.Error("proto unmarshal error %v", err)
		return
	}
	logs.Debug("HelloRouter message is %s", hello.Hello)

	err = req.GetConn().SendMsg(req.GetConn().ConnId(), nil)
}

func (this *HelloRouter) PostHandle(request iface.IRequest) {
	log.Debug("Call Router PostHandle")
}