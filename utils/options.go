package utils

import (
	"cginx/iface"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type serverOptions struct {
	TcpServer iface.Iserver
	Version   string

	Name           string `json:"name"`
	Host           string `json:"host"`
	TcpPort        int    `json:"tcp_port"`
	MaxConn        int    `json:"max_conn"`
	MaxPackageSize uint32 `json:"max_package_size"` //数据包最大值
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
	}

	ServerOpt.LoadOptions()
}
