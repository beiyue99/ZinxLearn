package znet

import (
	"fmt"
	"net"
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
	//当前server的router
	Router ziface.IRouter
}

func (s *Server) Start() {
	fmt.Printf("[Start] Server Listenner ar IP:%s,Port:%d is staring\n", s.IP, s.Port)
	go func() {
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
			//如果有客户端连接进来，accept会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			//将处理方法和conn绑定
			dealConn := NewConnection(conn, cid, s.Router)
			cid++
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {

	//将服务器的一些资源释放

}

func (s *Server) Serve() {
	//启动server的服务功能
	s.Start()
	// 阻塞主函数，可以做一些启动服务器之后的额外业务
	select {}
}

func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
	fmt.Println("Add router suss!")
}

// 初始化Server类的方法
func NewServer(name string) ziface.Iserver {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
		Router:    nil,
	}
	return s
}
