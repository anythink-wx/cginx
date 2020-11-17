package cnet

import (
	"cginx/iface"
	"fmt"
	"strconv"
)

//消息分发模块实现
type MsgHandle struct {
	//存放消息id对应的 irouter
	apiList map[uint16]iface.Irouter
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{apiList: make(map[uint16]iface.Irouter, 10)}
}

func (m *MsgHandle) DoMsgHandler(req iface.Irequest) {
	id := req.GetMsgId()
	handler, ok := m.apiList[id]
	if !ok {
		fmt.Println("msgID", id, "notfound")
		return
	}
	handler.Before(req)
	handler.Handler(req)
	handler.After(req)
}

func (m *MsgHandle) AddRouter(msgId uint16, irouter iface.Irouter) {

	_, ok := m.apiList[msgId]
	if ok {
		fmt.Println("[waring] 消息对应的路由已经被添加 msgID:" + strconv.Itoa(int(msgId)))
	} else {
		m.apiList[msgId] = irouter
		fmt.Println("[info] add router msgID:" + strconv.Itoa(int(msgId)))

	}
}

func (m *MsgHandle) Count() int {
	return len(m.apiList)
}
