package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)

//连接管理模块

type ConnManager struct {
	connections map[uint32]ziface.IConnection //所有连接的集合
	connLock    sync.RWMutex                  //保护连接集合的读写锁
}

//创建连接管理器实例

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

// 添加连接
func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	connMgr.connections[conn.GetConnID()] = conn
	fmt.Println("connID=", conn.GetConnID(), " add to ConnManager succss :conn num=", connMgr.Len())
}

// 删除链接
func (connMgr *ConnManager) Remove(conn ziface.IConnection) {

	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	delete(connMgr.connections, conn.GetConnID())
	fmt.Println("connID=", conn.GetConnID(), " Remove to ConnManager succss :conn num=", connMgr.Len())
}

// 根据ConnID获取连接
func (connMgr *ConnManager) Get(connID uint32) (ziface.IConnection, error) {

	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()
	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connction not found!")
	}
}

// 得到连接总数
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)

}

// 清除所有连接
func (connMgr *ConnManager) ClearConn() {

	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	for connID, conn := range connMgr.connections {
		conn.Stop()
		delete(connMgr.connections, connID)
	}
	fmt.Println("Clear all connction succ!,conn num is", connMgr.Len())

}
