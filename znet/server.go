package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

// iServer 接口实现，定义一个Server服务器模块
type Server struct {
	//服务器的名称
	Name string
	//服务绑定的ip版本 tcp4 or other
	IPVersion string
	//服务绑定的IP地址
	IP string
	//服务绑定的端口
	Port int

	//当前的Server添加一个router，server注册的连接对应的处理业务
	Router ziface.IRouter
}

//============== 实现 ziface.IServer 里的全部接口方法 ========

func (s *Server) Start() {
	fmt.Printf("[GZinx] Server Name : %s, listener at IP : %s, Port : %d is starting\n",
		utils.GlobalObject.Name,
		utils.GlobalObject.Host,
		utils.GlobalObject.TcpPort)
	fmt.Printf("[GZinx] Version %s, MaxConn : %d, MaxPackeetSize : %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPackageSize)

	go func() {
		//	1获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error : ", err)
			return
		}

		//	2尝试监听服务器的地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, " error : ", err)
			return
		}
		fmt.Println("start zinx server success, ", s.Name, " success, Listening ...")

		var cid uint32
		cid = 0

		//	3阻塞地等待客户端链接，处理客户端的链接业务（读写）
		for {
			//如果有客户端链接过来，阻塞会返回
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err： ", err)
				continue
			}
			// 将处理新连接的业务方法 和 conn 进行绑定 得到我们的连接模块
			dealConn := NewConnection(conn, cid, s.Router)
			cid++

			// 启动当前的连接业务处理
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server , name ", s.Name)
	//TODO 将一些服务器的资源、状态或者一些已经开辟的连接信息 进行停止或者回收
}

func (s *Server) Serve() {
	//启动服务功能
	s.Start()

	//TODO 座椅写启动服务器之后的额外业务

	//阻塞状态
	select {}
}

func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
	fmt.Println("Add Router Success!")
}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		Router:    nil,
	}

	return s
}
