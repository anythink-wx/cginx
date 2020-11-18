package iface

type Iserver interface {
	Start()
	Stop()
	Serve()

	AddRouter(msgId uint16, router Irouter)

	GetConnMgr() IconnManager

	//OnOpen(iconnection Iconnection)
	//OnClose()

	SetConnOpen(func(conn Iconnection))
	SetConnClose(func(conn Iconnection))
	CallConnOpen(conn Iconnection)
	CallConnClose(conn Iconnection)
}
