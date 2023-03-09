package znet

import "zinx/ziface"

type Request struct {
	//已经和客户建立好的连接
	conn ziface.IConnection

	//客户端请求的数据
	data []byte
}

// 得到当前连接
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// 得到请求的消息数据
func (r *Request) GetData() []byte {
	return r.data
}
