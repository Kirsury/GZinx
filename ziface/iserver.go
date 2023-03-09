package ziface

// 定义服务器接口
type IServer interface {
	//启动服务器方法
	Start()
	//停止服务器方法
	Stop()
	//运行服务器方法
	Serve()

	//路由功能：给当前的服务注册一个路由方法，供客户端的连接处理使用
	AddRouter(router IRouter)
}
