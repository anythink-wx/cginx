package iface

type Iserver interface {
	Start()
	Stop()
	Serve()

	AddRouter(router Irouter)
}
