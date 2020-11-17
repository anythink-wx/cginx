package iface

type Iserver interface {
	Start()
	Stop()
	Serve()

	AddRouter(msgId uint16, router Irouter)
}
