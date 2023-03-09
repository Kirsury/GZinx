package ziface

/*
 IRequest 接口：
实际上是吧客户端请求的连接信息，和 请求的数据包装到一个请求中
*/

type IRequest interface {
	//得到当前连接
	GetConnection() IConnection

	//得到请求的消息数据
	GetData() []byte
}
