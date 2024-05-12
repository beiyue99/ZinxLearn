package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

type Connection struct {
	//当前conn属于哪个server
	TcpServer ziface.Iserver
	//当前连接tcp连接
	Conn *net.TCPConn
	//当前连接的ID
	ConnID uint32
	//当前连接状态
	isClosed bool
	//通讯的channel,用于告知写groutine连接已经关闭
	ExitChan chan bool
	//用于读写Groutine之间的通信
	msgChan chan []byte
	//管理消息id的对应处理方法router
	MsgHandle ziface.IMsgHandle
}

func NewConnection(server ziface.Iserver, conn *net.TCPConn, connID uint32, msghandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer: server,
		Conn:      conn,
		ConnID:    connID,
		MsgHandle: msghandler,
		isClosed:  false,
		ExitChan:  make(chan bool, 1),
		msgChan:   make(chan []byte),
	}
	c.TcpServer.GetConnMgr().Add(c)
	return c
}

// 连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader groutine is runing...")
	defer fmt.Println("connID", c.ConnID, "Reader is exit,RemoteAddr is", c.RemoteAddr().String())
	defer c.Stop()
	for {
		//创建一个拆包解包对象
		dp := NewDataPack()

		//读取客户端的MsgHead
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error")
			break
		}
		//拆包，得到len和id
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("Unpack err", err)
			break
		}
		//根据datalen,读取数据内容
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error", err)
				break
			}
		}
		msg.SetData(data)

		//得到当前conn的请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}
		if utils.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHandle.SendMsgToTaskQueue(&req)
		} else {

			//从路由中，找到对应的router调用
			go c.MsgHandle.DoMsgHandler(&req)
		}
	}
}

//提供一个SendMsg方法，待发送的数据先封包，再发送

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection closed when send msg")
	}
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id=", msgId)
		return errors.New("pack error msg")
	}
	//将数据发送给管道
	c.msgChan <- binaryMsg
	return nil
}

// 写消息的groutine，专门把消息发给客户端
func (c *Connection) StartWriter() {
	fmt.Println("Write groutine is running!")
	defer fmt.Println("[conn Write groutine exit!]", c.RemoteAddr().String())

	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data err", err)
				return
			}
		case <-c.ExitChan:
			return
		}
	}
}

// 启动连接
func (c *Connection) Start() {
	fmt.Println("Conn start...ConnID:", c.ConnID)
	go c.StartReader()
	go c.StartWriter()
}

// 停止连接
func (c *Connection) Stop() {
	if c.isClosed {
		return
	}
	fmt.Println("Conn stop,ConnID=", c.ConnID)
	c.isClosed = true
	c.Conn.Close()
	c.ExitChan <- true
	c.TcpServer.GetConnMgr().Remove(c)
	close(c.ExitChan)
	close(c.msgChan)

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
