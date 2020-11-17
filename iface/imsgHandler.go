package iface
//消息管理抽象层




type ImsgHandle interface {
	//执行对应router消息

	DoMsgHandler(req Irequest)

	//为消息添加处理的路由

	AddRouter(msgId uint16, irouter Irouter )

	//路由数量
	Count() int
}
