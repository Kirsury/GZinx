package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

/*
基于zinx框架来开发的服务端应用程序
*/

// ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// Test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	//读取客户端的数据，再回写
	fmt.Println("recv from client: msgId = ", request.GetMsgID(), ", data = ", string(request.GetData()))
	if err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping")); err != nil {
		fmt.Println(err)
	}
}

// Hello GZinx test 自定义路由
type HelloGZinxRouter struct {
	znet.BaseRouter
}

// Test Handle
func (this *HelloGZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloGZinxRouter Handle...")
	//读取客户端的数据，再回写
	fmt.Println("recv from client: msgId = ", request.GetMsgID(), ", data = ", string(request.GetData()))
	if err := request.GetConnection().SendMsg(201, []byte("Helo Welcome to Gzinx!!")); err != nil {
		fmt.Println(err)
	}
}

// 创建连接之后执行钩子函数
func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("---> DoConnectionBegin is Called ...")
	if err := conn.SendMsg(202, []byte("DoConnectionBegin BEGIN")); err != nil {
		fmt.Println(err)
	}

	//给当前的连接设置一些属性
	fmt.Println("Set conn property ...")
	conn.SetProperty("Name", "jojo-大乔")
	conn.SetProperty("Github", "https://github.com/Kirsury/GZinx")
	conn.SetProperty("Home", "www.bing.cn")
}

// 连接断开之前的需要执行的函数
func DoConnectionLost(conn ziface.IConnection) {
	fmt.Println("---> DoConnectionLost is Called")
	fmt.Println("conn ID = ", conn.GetConnID(), " is lost")

	//获取连接属性
	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Name = ", name)
	}
	if github, err := conn.GetProperty("Github"); err == nil {
		fmt.Println("Github = ", github)
	}
	if home, err := conn.GetProperty("Home"); err == nil {
		fmt.Println("Home = ", home)
	}

}

func main() {
	//1 创建一个server句柄，使用Zinx的api
	s := znet.NewServer("[zinx V0.10]")

	//2 注册连接Hook钩子函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	//3 给当前server添加一个自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloGZinxRouter{})

	//4 启动server
	s.Serve()
}
