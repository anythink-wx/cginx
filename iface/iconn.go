package iface

import "net"

type Iconnection interface {

	//启动连接
	Open(requestId *uint64)

	//停止连接
	Close()

	//获取连接的 socket 套接字
	GetTCPConnection() *net.TCPConn

	//获取连接id
	GetConnID() uint32

	//获取远程客户端tcp状态
	RemoteAddr() net.Addr

	//发送数据
	Send(data []byte) error

	SendMsg( uint16,  []byte) error

	SetProp(k string, val interface{})
	GetProp(k string) (interface{})
	DelProp(k string)
}


//处理业务的类型
//type HandleFunc func(*net.TCPConn, []byte, int) error
