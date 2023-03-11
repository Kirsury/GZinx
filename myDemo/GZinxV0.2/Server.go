package main

import "zinx/znet"

/*
基于zinx框架来开发的服务端应用程序
*/

func main() {
	//	1创建一个server句柄，使用Zinx的api
	s := znet.NewServer("[zinx V0.2]")
	//	2启动server
	s.Serve()
}
