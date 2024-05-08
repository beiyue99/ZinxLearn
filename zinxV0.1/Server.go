package main

import "zinx/znet"

//基于zinx框架开发的服务器应用程序

func main() {
	//1 创建一个server,使用zinx的api
	s := znet.NewServer("[zinx V0.1]")
	//2 启动server
	s.Serve()
}
