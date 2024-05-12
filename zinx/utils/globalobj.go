package utils

import (
	"encoding/json"
	"os"
	"zinx/ziface"
)

//存储全局参数，供其他模块使用，由用户配置

type GlobalObj struct {
	TcpServer        ziface.Iserver
	Host             string //监听的IP
	TcpPort          int    //端口
	Name             string //服务器名称
	Version          string //版本号
	MaxConn          int    //最大连接数
	MaxPackageSize   uint32 //数据包最大值
	WorkerPoolSize   uint32 //当前业务工作Worker池的groutine数量
	MaxWorkerTaskLen uint32 //允许开辟的最大groutine数量
}

//定义一个全局的对外Globalobj

var GlobalObject *GlobalObj

// 从zinx.json加载自定义的参数
func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	//将json文件解析到struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

/*
提供一个init方法，初始化当前的GlobalObject

*/

func init() {
	//如果配置文件没有加载，就是这样的默认值
	GlobalObject = &GlobalObj{
		Name:             "ZinxServerApp",
		Version:          "V0.9",
		TcpPort:          8999,
		Host:             "0.0.0.0",
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
	}
	GlobalObject.Reload()
}
