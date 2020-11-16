package iface

//路由模块
type Irouter interface {
	Before(r Irequest)
	Handler(r Irequest)
	After(r Irequest)
}
