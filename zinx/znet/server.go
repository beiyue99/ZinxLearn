package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

//Iserver的接口实现，定义一个Server的服务器类

type Server struct {
	//服务器名称
	Name string
	//服务器绑定的ip版本
	IPVersion string
	//服务器监听的ip
	IP string
	//服务器监听的端口
	Port int
	//当前sever的消息管理模块，用来绑定MsgId和处理业务API的关系
	MsgHandler ziface.IMsgHandle
	//Server的连接管理器
	ConnMgr ziface.IConnManager
	//创建连接之后的hook函数
	OnConnStart func(conn ziface.IConnection)
	//连接销毁之前的hook函数
	OnConnStop func(conn ziface.IConnection)
}

func (s *Server) Start() {

	fmt.Printf("[zinx] Server Name:%s,Listen at %s,%d\n", utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[zinx] Version is %s,Maxconn is %d,MaxPackageSize is %d\n", utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)

	go func() {
		//开启Worker工作池
		s.MsgHandler.StartWorkerPool()

		//获取一个TCP的地址
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error:", err)
			return
		}
		//监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
			return
		}
		fmt.Println("start Zinx server succ,", s.Name, "succ,Listenning...")
		cid := uint32(0)
		//阻塞等待客户端连接，处理客户端业务
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				fmt.Println("===>too many conn,connnum is", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}

			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++
			go dealConn.Start()
		}
	}()
}

// 调用连接管理器的ClearConn释放资源
func (s *Server) Stop() {

	//将服务器的一些资源释放
	fmt.Println("[stop] Zinx server name", s.Name)
	s.ConnMgr.ClearConn()
}

// 调用s.Start启动server的服务功能
func (s *Server) Serve() {
	s.Start()
	// 阻塞主函数，可以做一些启动服务器之后的额外业务
	select {}
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add router suss!")
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

// 初始化Server类的方法
func NewServer(name string) ziface.Iserver {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
	return s
}

// 注册hook方法
func (s *Server) SetOnConnStart(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStart = hookFunc
}
func (s *Server) SetOnConnStop(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// 调用hook方法
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---->call OnConnStart!")
		s.OnConnStart(conn)
	}
}
func (s *Server) CallOnConnStop(conn ziface.IConnection) {

	if s.OnConnStop != nil {
		fmt.Println("---->call OnConnStop!")
		s.OnConnStop(conn)
	}
}
