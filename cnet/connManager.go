package cnet

import (
	"cginx/iface"
	"errors"
	"fmt"
	"strconv"
	"sync"
)

type ConnManager struct {
	connections map[uint32]iface.Iconnection
	connLock    sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]iface.Iconnection),
	}
}

func (c *ConnManager) Remove(conn iface.Iconnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()
	delete(c.connections, conn.GetConnID())
	fmt.Println("Remove conn to manager", c)
}

func (c *ConnManager) Add(conn iface.Iconnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	c.connections[conn.GetConnID()] = conn
	fmt.Println("add conn to manager", c)
}

func (c *ConnManager) Get(connID uint32) (conn iface.Iconnection, err error) {
	c.connLock.RLock()
	defer c.connLock.RUnlock()

	conn, ok := c.connections[connID]
	if !ok {
		err = errors.New("conn not found" + strconv.Itoa(int(connID)))
		return nil, err
	}

	return
}

func (c *ConnManager) Len() int {
	//c.connLock.RLock()
	//defer c.connLock.RUnlock()

	return len(c.connections)+1
}

func (c *ConnManager) ClearConn() {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	for connId, conn := range c.connections {
		conn.Close()
		delete(c.connections, connId)
	}
	fmt.Println("[connManager]close all connections")
}
