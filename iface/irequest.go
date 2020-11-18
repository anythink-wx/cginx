package iface

//request 请求封装， 将连接和数据绑定在一起
type Irequest interface {
	//连接句柄

	//请求的内容

	//方法 获取连接  获取数据

	GetConnection() Iconnection

	//需要进行消息分装
	GetData() []byte

	GetMsgId() uint16

	//获取request序号
	GetReqID() uint64

}
