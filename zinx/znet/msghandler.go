package znet

import (
	"zinx/log"
	"zinx/ziface"
)

type MsgHandle struct {
	Apis map[uint32] ziface.IRouter	// 存放每个MsgId 所对应的处理方法的map属性
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

// 马上以非阻塞方式处理消息
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		log.Error("api msgId %d is not fount", request.GetMsgID())
		return
	}

	// 执行对应处理方法
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {
	// 1. 判断当前msg绑定的api处理方法是否已经存在
	if _, ok := mh.Apis[msgId]; ok {
		panic("repeated api, msgid is " + string(msgId))
	}
	// 2.添加msg与api的绑定关系
	mh.Apis[msgId] = router
	log.Release("Add api msgid = %d", msgId)
}