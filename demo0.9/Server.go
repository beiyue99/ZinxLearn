package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

// ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Ping router Handle...")
	//先打印读取到的客户端的数据，再给客户端回  ping ping ping
	fmt.Println("recv form client :msgId=", request.GetMsgID(), "data=", string(request.GetData()))
	err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println("err")
	}
}

// Hello test 自定义路由
type HelloRouter struct {
	znet.BaseRouter
}

func (this *HelloRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Hello router Handle...")
	//先打印读取到的客户端的数据，再给客户端回
	fmt.Println("recv form client :msgId=", request.GetMsgID(), "data=", string(request.GetData()))
	err := request.GetConnection().SendMsg(201, []byte("Hello welcome to zinx..."))
	if err != nil {
		fmt.Println("err")
	}
}

func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("=====>DoConnectionBegin is called!")
	if err := conn.SendMsg(202, []byte("DoConnection Begin")); err != nil {
		fmt.Println(err)
	}
}
func DoConnectionLost(conn ziface.IConnection) {
	fmt.Println("=====>DoConnectionLost is called!")
	fmt.Println("conn ID is", conn.GetConnID(), "lost...")
}

func main() {
	// 创建一个server,使用zinx的api
	s := znet.NewServer("[zinx V0.6]")
	//注册hook函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)
	//添加一个router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})
	// 启动server
	s.Serve()
}
