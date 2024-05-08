package znet

import (
	"fmt"
	"net"
	"zinx/ziface"
)

type Connection struct {
	//当前连接tcp连接
	Conn *net.TCPConn
	//当前连接的ID
	ConnID uint32
	//当前连接状态
	isClosed bool
	//通讯的channel
	ExitChan chan bool
	//连接的处理方法 router
	Router ziface.IRouter
}

func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		Router:   router,
		isClosed: false,
		ExitChan: make(chan bool, 1),
	}
	return c
}

// 连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader groutine is runing...")
	defer fmt.Println("connID", c.ConnID, "Reader is exit,RemoteAddr is", c.RemoteAddr().String())
	defer c.Stop()
	for {
		buf := make([]byte, 512)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("Read err:", err)
			continue
		}
		//得到当前conn的请求数据
		req := Request{
			conn: c,
			data: buf,
		}
		//从路由中，找到对应的router调用
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
	}
}

// 启动连接
func (c *Connection) Start() {
	fmt.Println("Conn start...ConnID:", c.ConnID)
	go c.StartReader()
}

// 停止连接
func (c *Connection) Stop() {
	if c.isClosed {
		return
	}
	fmt.Println("Conn stop,ConnID=", c.ConnID)
	c.isClosed = true
	c.Conn.Close()
	close(c.ExitChan)

}

// 获取连接
func (c *Connection) GetTCPConnection() *net.TCPConn {

	return c.Conn
}

// 获取连接ID
func (c *Connection) GetConnID() uint32 {

	return c.ConnID
}

// 获取地址和端口
func (c *Connection) RemoteAddr() net.Addr {

	return c.Conn.RemoteAddr()
}

// 发送数据
func (c *Connection) Sent(data []byte) error {
	return nil
}
