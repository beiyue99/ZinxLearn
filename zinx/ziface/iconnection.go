package ziface

import "net"

type IConnection interface {
	//启动连接
	Start()
	//停止连接
	Stop()
	//获取连接
	GetTCPConnection() *net.TCPConn
	//获取连接ID
	GetConnID() uint32
	//获取地址和端口
	RemoteAddr() net.Addr
	//发送数据
	SendMsg(msgId uint32, data []byte) error
}

type HandFunc func(*net.TCPConn, []byte, int) error
