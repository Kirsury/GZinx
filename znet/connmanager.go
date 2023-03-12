package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)

/*
连接管理模块
*/
type ConnManager struct {
	connections map[uint32]ziface.IConnection //管理的连接集合
	connLock    sync.RWMutex                  //保护连接集合的读写锁
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

// 添加连接
func (cm *ConnManager) Add(conn ziface.IConnection) {
	//保护共享资源map，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	//将conn加入到ConnManager中
	cm.connections[conn.GetConnID()] = conn
	fmt.Println("connID = ", conn.GetConnID(), " add to ConnManager successfully: conn num = ", cm.Len())
}

// 删除连接
func (cm *ConnManager) Remove(conn ziface.IConnection) {
	//保护共享资源map，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	//删除连接信息
	delete(cm.connections, conn.GetConnID())
	fmt.Println("connID = ", conn.GetConnID(), " remove from ConnManager successfully: conn num = ", cm.Len())
}

// 根据连接ID查找对应的连接
func (cm *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	//保护共享资源map，加读锁
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()

	if conn, ok := cm.connections[connID]; ok {
		//找到了
		return conn, nil
	} else {
		return nil, errors.New("connection not FOUND")
	}
}

// 总连接个数
func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

// 清理终止全部的连接
func (cm *ConnManager) ClearConn() {
	//保护共享资源map，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	//删除conn并停止conn的工作
	for connID, conn := range cm.connections {
		//停止
		conn.Stop()
		//删除
		delete(cm.connections, connID)
	}

	fmt.Println("Clear All connections succ! conn num = ", cm.Len())
}
