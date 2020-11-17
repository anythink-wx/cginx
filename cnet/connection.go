package cnet

import (
	"cginx/iface"
	"errors"
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

func (c *Connection) SendMsg(id uint16, data []byte) (err error) {
	if c.isClose {
		return errors.New("conn is closed")
	}

	m := NewMessagePackage(id, data)
	byteMsg, err := NewDataPack().Pack(m)
	if err != nil {
		return
	}
	_, err = c.GetTCPConnection().Write(byteMsg)
	if err != nil {
		fmt.Println("send msg", m.GetMsgId(), "error")
		return
	}
	return
}

//读取客户消息
func (c *Connection) readerGoroutine() {
	fmt.Println("readerGoroutine is running ... ")
	defer fmt.Println("readerGoroutine is exit connID=", c.ConnID)
	defer c.Close()

	for {

		//拆包
		pack := NewDataPack()
		//读取header 6 字节
		buf := make([]byte, pack.GetHeadLen())
		_, err := c.GetTCPConnection().Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("readerGoroutine error", err)
			break
		}
		msgHead, err := pack.Unpack(buf)
		if err != nil {
			fmt.Println("Unpack header error", err)
			break
		}

		msg := msgHead.(*message)

		//GetMsgLen > 0 说明msg有数据
		if msgHead.GetMsgLen() > 0 {
			msg.Data = make([]byte, msgHead.GetMsgLen())

			if _, err = c.GetTCPConnection().Read(msg.Data); err != nil {
				fmt.Println("Unpack data error", err)
				break
			}
		}

		//得到Irequest 结构
		req := Request{
			conn: c,
			msg:  msg,
		}

		go func(r iface.Irequest) {
			c.Router.Before(r)
			c.Router.Handler(r)
			c.Router.After(r)
		}(&req)

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
