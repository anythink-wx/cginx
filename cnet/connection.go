package cnet

import (
	"cginx/iface"
	"fmt"
	"io"
	"net"
)

type Connection struct {
	Conn    *net.TCPConn
	ConnID  uint32
	isClose bool
	//处理业务的方法
	HandleCall iface.HandleFunc
	//是否关闭
	ExitChan chan bool
}

//构造方法
func NewConnection(conn *net.TCPConn, connId uint32, callback iface.HandleFunc ) *Connection {

	return &Connection{
		Conn:       conn,
		ConnID:     connId,
		HandleCall: callback,
		ExitChan:   make(chan bool),
	}

}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) RemoteAddr() net.Addr {

	return c.Conn.RemoteAddr()
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) Send(data []byte) error {

	return nil
}

//读取客户消息
func (c *Connection) readerGoroutine() {
	fmt.Println("readerGoroutine is running ... ")
	defer fmt.Println("readerGoroutine is exit connID=", c.ConnID)
	defer c.Close()

	for {

		buf := make([]byte, 512)
		ln, err := c.Conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("readerGoroutine error", err)
			continue
		}


		err = c.HandleCall(c.Conn, buf, ln)
		if err != nil {
			fmt.Println("HandleCall error", err)
			break
		}

	}

}
func (c *Connection) open() {
	fmt.Println("conn open  connID=", c.ConnID)
	// 启动 读取客户数据的 go程
	go c.readerGoroutine()

	//启动会写数据的go程

}

func (c *Connection) Close() {

	fmt.Println("conn Stop, connID=", c.ConnID)
	if c.isClose == true {
		return
	}
	c.isClose = true
	c.Conn.Close()
	close(c.ExitChan)

}
