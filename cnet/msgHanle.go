package cnet

import (
	"cginx/iface"
	"cginx/utils"
	"fmt"
	"strconv"
)

//消息分发模块实现
type MsgHandle struct {
	//存放消息id对应的 irouter
	apiList map[uint16]iface.Irouter
	//worker 工作池的数量
	WorkerPoolSize uint8

	TaskQueue []chan iface.Irequest
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		apiList:        make(map[uint16]iface.Irouter),
		WorkerPoolSize: utils.ServerOpt.WorkerPoolSize,
		TaskQueue:      make([]chan iface.Irequest, utils.ServerOpt.WorkerPoolSize),
	}
}

//开启worker 工作池
func (m *MsgHandle) StartWorkerPool() {
	for i := 0; i < int(utils.ServerOpt.WorkerPoolSize); i++ {
		//给管道开辟内存空间， 不然空指针了
		m.TaskQueue[i] = make(chan iface.Irequest, utils.ServerOpt.MaxPackageSize)
		go m.startWorker(i)
	}
}

func (m *MsgHandle) startWorker(queueId int) {
	iRequestChan := m.TaskQueue[queueId]
	fmt.Println("start business worker", queueId)

	//不停的阻塞等待消息
	for {
		select {
		case req := <-iRequestChan:
			m.DoMsgHandler(req)
		}
	}

}


func (m *MsgHandle) PushWorkerQueue(req iface.Irequest) {

	//将消息平均分配
	size := req.GetReqID() % uint64(utils.ServerOpt.WorkerPoolSize)
	// 将该消息发送对应worker的 TaskQueue
	//fmt.Println("push msg to WorkerQueue ConnID=",req.GetConnection().GetConnID()," ReqID:",req.GetReqID(),"workerID:",size)

	irequestsChan := m.TaskQueue[size]
	irequestsChan<-req
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
