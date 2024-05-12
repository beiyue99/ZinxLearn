package ziface

// 消息管理抽象层

type IMsgHandle interface {
	//执行对应的消息处理方法
	DoMsgHandler(request IRequest)
	//为消息添加路由
	AddRouter(msgID uint32, router IRouter)
	//启动Worker工作池
	StartWorkerPool()
	//将消息发给消息队列处理
	SendMsgToTaskQueue(request IRequest)
}
