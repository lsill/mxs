package mnet

import (
	"mxs/gamex/api/websocket/iface"
	"mxs/gamex/utils"
	logs "mxs/log"
)

type MsgHandle struct {
	Apis map[int32]iface.IRouter      // 存放每个MsgId 所对应的处理方法的map属性
	WorkerPoolSize uint32              // 业务工作worker池的数量
	TaskQueue	[]chan iface.IRequest // Worker负责取任务的消息队列
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[int32]iface.IRouter),
		WorkerPoolSize: utils.GloUtil.MaxWorkerTaskLen,
		// 一个worker对应一个queue
		TaskQueue: make([]chan iface.IRequest, utils.GloUtil.MaxWorkerTaskLen),
	}
}

func (mh *MsgHandle) DoMsgHandler(req iface.IRequest) {
	handler, ok := mh.Apis[req.GetPkTyp()]
	if !ok {
		logs.Error("api msgid %d is not fount", req.GetPkTyp())
		return
	}
	// 执行对应的处理方法
	handler.PreHandle(req)
	handler.Handle(req)
	handler.PostHandle(req)
}

// 为消息添加对应的处理逻辑
func (mh *MsgHandle) AddRouter(typ int32, router iface.IRouter) {
	if _, ok := mh.Apis[typ]; ok {
		logs.Error("repeated api, typ is %d", typ)
		return
	}
	mh.Apis[typ] = router
	logs.Release("add api typ = %d", typ)
}

// 启动一工作流程
func (mh *MsgHandle) StartOneWorker(workerid int, taskQueue chan iface.IRequest) {
	logs.Debug("workerid %d is Started", workerid)
	for {
		select {
		case req := <- taskQueue:
			mh.DoMsgHandler(req)
		}
	}
}

// 启动工作池
func(mh *MsgHandle) StarWorkerPool() {
	// 遍历需要启动worker的数量 依次启动
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		mh.TaskQueue[i] = make(chan iface.IRequest, 1000)
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

func (mh *MsgHandle) SendMsgToTaskQueue(req iface.IRequest) {
	workerid := req.GetConn().ConnId() % mh.WorkerPoolSize
	logs.Debug("add connid %d req msgtyp is %v", req.GetConn().ConnId(), req.GetPkTyp())
	mh.TaskQueue[workerid] <- req
}

