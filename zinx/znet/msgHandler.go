package znet

import (
	"fmt"
	"strconv"
	"zinx/ziface"
)

//消息处理模块的实现

type MsgHandle struct {
	Apis map[uint32]ziface.IRouter
}

// 初始化方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
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
