package znet

import (
	"fmt"
	"strconv"
	"zinx/utils"
	"zinx/ziface"
)

//消息处理模块的实现

type MsgHandle struct {
	//存放每个msgId和对应的处理方法的map
	Apis map[uint32]ziface.IRouter
	//负责存放消息的消息队列
	TaskQueue []chan ziface.IRequest
	//work数量
	WorkerPoolSize uint32
}

// 初始化方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

// 执行对应的消息处理方法
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	//先找到msgId
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api MsgId=", request.GetMsgID(), "is not found!")
		return
	}
	//根据MsgID调用对应的处理API
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 为消息添加路由
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	//判断当前ID的处理API是否已经存在
	if _, ok := mh.Apis[msgID]; ok {
		panic("repeat api,msgId=" + strconv.Itoa(int(msgID)))
	}
	//添加ID和API的绑定关系
	mh.Apis[msgID] = router
	fmt.Println("Add api MsgId=", msgID, "succ!")
}

//启动worker工作池

func (mh *MsgHandle) StartWorkerPool() {
	//根据WorkerPoolSize,开启那么多个worker
	for i := 0; i < int(mh.WorkerPoolSize); i++ {

		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}

}

// 启动一个Worker工作流程
// 通道读取操作 (for request := range taskQueue) 会阻塞，直到通道关闭或者有新的数据可供读取。这是 Go 通道的特性
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("WorkerId = ", workerID, "is started!")
	for request := range taskQueue {
		mh.DoMsgHandler(request)
	}

}

// 将消息发给TaskQueue
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {

	//根据ConnID来分配任务
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnId为", request.GetConnection().GetConnID(),
		"的消息Id为", request.GetMsgID(), "的消息 to WorkerID为", workerID, "的任务队列")
	//将消息发给worker的TaskQueue
	mh.TaskQueue[workerID] <- request
}
