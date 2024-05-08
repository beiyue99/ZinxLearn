package znet

import "zinx/ziface"

//先定义一个BaseRouter基类，然后根据需要重写基类方法

type BaseRouter struct{}

// 如果有的router没有pre和post业务，那么就无需重写
func (br *BaseRouter) PreHandle(request ziface.IRequest) {}

func (br *BaseRouter) Handle(request ziface.IRequest) {}

func (br *BaseRouter) PostHandle(request ziface.IRequest) {}
