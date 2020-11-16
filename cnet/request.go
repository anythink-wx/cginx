package cnet

import "cginx/iface"

type Request struct {
	//连接

	conn iface.Iconnection
	//客户端数据

	data []byte
}

func (r *Request) GetConnection() iface.Iconnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.data
}
