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
	if err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping\n")); err != nil {
		fmt.Println(err)
	}
}

func main() {
	//1 创建一个server句柄，使用Zinx的api
	s := znet.NewServer("[zinx V0.5]")

	//2 给当前server添加一个自定义的router
	s.AddRouter(&PingRouter{})

	s.Serve()
}
