package znet

import (
	"fmt"
	"net"
	"zinx/ziface"
)

type Connection struct {
	//当前连接tcp套间字
	Conn *net.TCPConn
	//当前连接的ID
	ConnID uint32
	//当前连接状态
	isClosed bool
	//回调函数
	handleAPI ziface.HandFunc
	//通讯的channel
	ExitChan chan bool
}

func NewConnection(conn *net.TCPConn, connID uint32, callback_api ziface.HandFunc) *Connection {
	c := &Connection{
		Conn:      conn,
		ConnID:    connID,
		handleAPI: callback_api,
		isClosed:  false,
		ExitChan:  make(chan bool, 1),
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
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("Read err:", err)
			continue
		}
		//读完数据，调用绑定的handleAPI处理
		if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
			fmt.Println("ConnID", c.ConnID, "handle is err", err)
			break
		}
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
