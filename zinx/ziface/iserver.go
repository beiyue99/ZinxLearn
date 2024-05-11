package ziface

// 定义一个服务器接口
type Iserver interface {

	//启动服务器
	Start()
	//停止服务器
	Stop()
	//运行服务器
	Serve()
	//路由器
	AddRouter(msgID uint32, router IRouter)
}
