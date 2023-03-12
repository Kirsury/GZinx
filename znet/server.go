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

	//当前Server的消息管理模块，用来绑定MsgID和对应的处理业务的API关系
	MsgHandler ziface.IMsgHandle

	//该server的连接管理器
	ConnMgr ziface.IConnManager

	//该server创建连接之后自动调用的Hook函数 -- OnConnStart
	OnConnStart func(conn ziface.IConnection)
	//该server销毁连接之前自动调用的Hook函数 -- OnConnStop
	OnConnStop func(conn ziface.IConnection)
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
		//0 开启消息队列及worker工作池
		s.MsgHandler.StartWorkerPool()

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

			//设置最大连接个数判断，如果超过最大连接，那么则关闭此新的连接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				//TODO 给客户端响应一个超出最大连接的错误包
				fmt.Println("Too Many Connections MaxConn = ", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}

			// 将处理新连接的业务方法 和 conn 进行绑定 得到我们的连接模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++

			// 启动当前的连接业务处理
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[STOP] GZinx server , name ", s.Name)
	//将一些服务器的资源、状态或者一些已经开辟的连接信息 进行停止或者回收
	s.ConnMgr.ClearConn()
}

func (s *Server) Serve() {
	//启动服务功能
	s.Start()

	//TODO 座椅写启动服务器之后的额外业务

	//阻塞状态
	select {}
}

func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
	fmt.Println("Add Router Success!")
}

func NewServer(name string) ziface.IServer {
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

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

// 注册OnConnStart钩子函数的方法
func (s *Server) SetOnConnStart(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// 注册OnConnStop钩子函数的方法
func (s *Server) SetOnConnStop(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// 调用OnConnStart钩子函数的方法
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("----> Call OnConnStart() ...")
		s.OnConnStart(conn)
	}
}

// 调用OnConnStop钩子函数的方法
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("----> Call OnConnStop() ...")
		s.OnConnStop(conn)
	}

}
