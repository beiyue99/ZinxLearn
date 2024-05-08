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

func (this *PingRouter) PreHandle(request ziface.IRequest) {

	fmt.Println("Call router PreHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("brefor ping...\n"))
	if err != nil {
		fmt.Println("call back brefor ping err")
	}
}

func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call router Handle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping ping ping...\n"))
	if err != nil {
		fmt.Println("call back ping ping  ping err")
	}

}

func (this *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call router PostHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("afrer ping...\n"))
	if err != nil {
		fmt.Println("call back after ping err")
	}

}

func main() {
	// 创建一个server,使用zinx的api
	s := znet.NewServer("[zinx V0.3]")
	//添加一个router
	s.AddRouter(&PingRouter{})
	// 启动server
	s.Serve()
}
