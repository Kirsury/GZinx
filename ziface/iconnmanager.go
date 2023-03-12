package ziface

/*
连接管理模块抽象层
*/

type IConnManager interface {
	//添加连接
	Add(conn IConnection)
	//删除连接
	Remove(conn IConnection)
	//根据连接ID查找对应的连接
	Get(connID uint32) (IConnection, error)
	//总连接个数
	Len() int
	//清理终止全部的连接
	ClearConn()
}
