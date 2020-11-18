package utils

import (
	"cginx/iface"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

type serverOptions struct {
	TcpServer iface.Iserver
	Version   string

	Name           string `json:"name"`
	Host           string `json:"host"`
	TcpPort        int    `json:"tcp_port"`
	MaxConn        int    `json:"max_conn"`
	MaxPackageSize uint32 `json:"max_package_size"` //数据包最大值

	WorkerPoolSize uint8 `json:"worker_pool_size"` //处理业务的worker池数量
	MaxWorkerTaskSize uint16  `json:"max_worker_task_size"` //每个worker的等待队列 有点 back_log 的意思？
}


func (s *serverOptions) LoadOptions (){

	root := filepath.Dir(os.Args[0])
	file, err := ioutil.ReadFile(root + "/cginx.json")
	if err != nil {
		fmt.Println("cginx.json not found, use default", err)
		return
	}

	err = json.Unmarshal(file, &ServerOpt)
	if err != nil {
		panic(err)
	}

}

var ServerOpt *serverOptions

func init() {
	ServerOpt = &serverOptions{
		Name:           "Cginx",
		Version:        "1.0.0",
		Host:           "0.0.0.0",
		TcpPort:        8810,
		MaxConn:        10240,
		MaxPackageSize: 4096,
		WorkerPoolSize:uint8(runtime.NumCPU()),
	}

	ServerOpt.LoadOptions()
}
