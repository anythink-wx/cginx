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

	//是否关闭
	ExitChan chan bool

	//当前该链接绑定的路由
	Router iface.Irouter
}

//构造方法
func NewConnection(conn *net.TCPConn, connId uint32, router iface.Irouter) *Connection {

	return &Connection{
		Conn:     conn,
		ConnID:   connId,
		Router:   router,
		ExitChan: make(chan bool),
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

		//得到Irequest 结构
		req := Request{
			conn: c,
			data: buf[:ln],
		}

		go func(r iface.Irequest) {
			c.Router.Before(r)
			c.Router.Handler(r)
			c.Router.After(r)
		}(&req)

		//err = c.HandleCall(c.Conn, buf, ln)
		//if err != nil {
		//	fmt.Println("HandleCall error", err)
		//	break
		//}

	}

}
func (c *Connection) Open() {
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
