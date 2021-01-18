package iface

type IMsgHandle interface {
	DoMsgHandler(requset IRequest)
	AddRouter(msgId int32, router IRouter) // 为消息添加具体的处理逻辑
	StarWorkerPool()                       // 启动worker工作池
	SendMsgToTaskQueue(request IRequest)   // 将消息交给TaskQueue,又worker进行处理
}

