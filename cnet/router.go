package cnet

import "cginx/iface"

//实现router 时嵌入这个base基类， 然后根据需要对该方法重写
type BaseRouter struct {
}

func (b *BaseRouter) Before(r iface.Irequest) {}
func (b *BaseRouter) Handler(r iface.Irequest)    {}
func (b *BaseRouter) After(r iface.Irequest)  {}
