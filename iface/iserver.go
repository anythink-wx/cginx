package iface

type Iserver interface {
	Start()
	Stop()
	Serve()
}
