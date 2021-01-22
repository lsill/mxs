package mnet

import (
	"mxs/log"
	"mxs/util"
	"mxs/util/api/kcp/iface"
)

/*
	WorkerPoolSize:TaskQueue中的每个队列应该是和一个Worker对应的，所以在创建TaskQueue中队列数量要和Worker的数量一致。
	TaskQueue:是一个Request请求信息的channel集合。用来缓冲提供worker调用的Request请求信息，worker会从对应的队列中获取客户端的请求数据并且处理掉。
 */
type MsgHandle struct {
	Apis map[uint32]iface.IRouter      // 存放每个MsgId 所对应的处理方法的map属性
	WorkerPoolSize uint32              // 业务工作worker池的数量
	TaskQueue	[]chan iface.IRequest // Worker负责取任务的消息队列
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]iface.IRouter),
		WorkerPoolSize: util.GloUtil.MaxWorkerTaskLen,
		// 一个worker对应一个queue
		TaskQueue: make([]chan iface.IRequest, util.GloUtil.MaxWorkerTaskLen),
	}
}

// 马上以非阻塞方式处理消息
func (mh *MsgHandle) DoMsgHandler(request iface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgTyp()]
	if !ok {
		log.Error("api msgtyp %d is not fount", request.GetMsgTyp())
		return
	}

	// 执行对应处理方法
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgId uint32, router iface.IRouter) {
	// 1. 判断当前msg绑定的api处理方法是否已经存在
	if _, ok := mh.Apis[msgId]; ok {
		panic("repeated api, msgid is " + string(msgId))
	}
	// 2.添加msg与api的绑定关系
	mh.Apis[msgId] = router
	log.Release("Add api msgid = %d", msgId)
}

// 启动一个worker工作流程
func(mh *MsgHandle) StarOneWorker(workerid int, taskQueue chan iface.IRequest) {
	log.Debug("workerid %d is started",workerid)
	for {
		select {
		// 有消息则取出队列的Request,并执行绑定方法
			case request := <-taskQueue:
				mh.DoMsgHandler(request)
		}
	}
}

// 启动worker工作池
func (mh *MsgHandle) StarWorkerPool() {
	// 遍历需要启动worker的数量，依次启动
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 一个worker被启动
		// 给当前worker对应的任务队列开辟空间
		mh.TaskQueue[i] = make(chan iface.IRequest, util.GloUtil.MaxWorkerTaskLen)
		// 启动当前worker,阻塞等待对应的任务队列是否有消息传递进来
		go mh.StarOneWorker(i, mh.TaskQueue[i])
	}
}

/*
		StartWorkerPool()方法是启动Worker工作池，这里根据用户配置好的WorkerPoolSize的数量来启动，
	然后分别给每个Worker分配一个TaskQueue，然后用一个goroutine来承载一个Worker的工作业务。
		StartOneWorker()方法就是一个Worker的工作业务，每个worker是不会退出的(目前没有设定
    worker的停止工作机制)，会永久的从对应的TaskQueue中等待消息，并处理。
 */

// 将消息交给taskQueue，由worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request iface.IRequest) {
	// 根据ConnID来分配当前的连接由哪个worker负责处理
	// 轮询的平均分配法则

	// 得到需要处理此条链接的workerid
	workerId := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	log.Debug("add Connid = %d request msgid = %d", request.GetConnection().GetConnID(), request.GetMsgTyp())
	mh.TaskQueue[workerId] <- request
}

/*
	SendMsgToTaskQueue()作为工作池的数据入口，这里面采用的是轮询的分配机制，
	因为不同链接信息都会调用这个入口，那么到底应该由哪个worker处理该链接的请求
	处理，整理用的是一个简单的求模运算。用余数和workerID的匹配来进行分配。
	最终将request请求数据发送给对应worker的TaskQueue，那么对应的worker的Goroutine就会处理该链接请求了。
 */