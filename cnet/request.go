package cnet

import "cginx/iface"

type Request struct {
	//连接

	conn iface.Iconnection
	//客户端数据

	//tlv消息包
	msg iface.Imessage

}

func (r *Request) GetConnection() iface.Iconnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetMsgData()
}

func (r *Request) GetMsgId() uint16 {
	return r.msg.GetMsgId()
}