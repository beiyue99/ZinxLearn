package ziface

//路由抽象接口

type IRouter interface {
	//处理conn业务之前的钩子方法

	PreHandle(request IRequest)

	//处理业务的主方法
	Handle(request IRequest)

	//处理conn业务之后的钩子方法
	PostHandle(request IRequest)
}
