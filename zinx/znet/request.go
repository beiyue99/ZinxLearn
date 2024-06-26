package znet

import "zinx/ziface"

type Request struct {
	//已经建立好的连接
	conn ziface.IConnection
	//客户端请求的数据
	msg ziface.IMessage
}

// 得到当前连接
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// 得到请求数据内容
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

//得到请求ID

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}
