package cnet

import (
	"cginx/iface"
	"cginx/utils"
	"errors"
	"fmt"
	"io"
	"net"
)

type Connection struct {
	TcpServer iface.Iserver
	Conn      *net.TCPConn
	ConnID    uint32
	isClose   bool

	//关闭连接
	ExitChan chan bool

	//创建无缓冲管道，用于读写分离
	msgChan chan []byte

	//当前该链接绑定的路由
	MsgHandle iface.ImsgHandle
}

//构造方法
func NewConnection(server iface.Iserver, conn *net.TCPConn, connId uint32, router iface.ImsgHandle) *Connection {

	c := &Connection{
		TcpServer: server,
		Conn:      conn,
		ConnID:    connId,
		MsgHandle: router,
		ExitChan:  make(chan bool),
		msgChan:   make(chan []byte),
	}

	c.TcpServer.GetConnMgr().Add(c)
	return c
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

	defer func() {
		i := recover()
		if i != nil {
			fmt.Println("SendMsg error=", i)
		}
	}()

	c.msgChan <- byteMsg

	return
}

//读取客户消息
func (c *Connection) readerGoroutine(requestId *uint64) {
	fmt.Println("readerGoroutine is running ... ")
	defer fmt.Println("readerGoroutine is exit connID=", c.ConnID)
	defer c.Close() //reader失败的时候释放该用户连接的chan资源

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
			ReqId: *requestId,
			conn:  c,
			msg:   msg,
		}
		*requestId++
		//将消息发送到工作池中 如果开启了工作池
		if utils.ServerOpt.WorkerPoolSize > 0 {
			go c.MsgHandle.PushWorkerQueue(&req)
		} else {
			go c.MsgHandle.DoMsgHandler(&req)
		}
	}

}

//创建一个只写的goroutine
func (c *Connection) writerGoroutine() {
	fmt.Println("writerGoroutine is running ")
	for {
		select {
		case msg := <-c.msgChan:
			_, err := c.GetTCPConnection().Write(msg)
			if err != nil {
				fmt.Println("send Data err:", err)
				return
			}

		case <-c.ExitChan:
			fmt.Println("revice c.ExitChan , writerGoroutine return")
			return
		}
	}
}

func (c *Connection) Open(RequestID *uint64) {

	fmt.Println("conn open  connID=", c.ConnID)
	// 启动 读取客户数据的 go程
	go c.readerGoroutine(RequestID)

	//启动会写数据的go程
	go c.writerGoroutine()

}

func (c *Connection) Close() {

	fmt.Println("conn Stop, connID=", c.ConnID)
	if c.isClose == true {
		return
	}
	c.isClose = true
	err := c.Conn.Close()
	if err != nil {
		fmt.Println("conn close error=", err)
	}
	close(c.ExitChan)
	close(c.msgChan)

	c.TcpServer.GetConnMgr().Remove(c)


}
