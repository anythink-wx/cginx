package cnet

import (
	"cginx/iface"
	"cginx/utils"
	"fmt"
	"net"
)

type Server struct {
	Name   string
	IPType string
	IP     string
	Port   int
	MsgHandle iface.ImsgHandle
}



func NewServer(name string, ) iface.Iserver{
	return &Server{
		Name:   utils.ServerOpt.Name,
		IPType: "tcp",
		IP:     utils.ServerOpt.Host,
		Port:   utils.ServerOpt.TcpPort,
		MsgHandle: NewMsgHandle(),
	}
}




func (s *Server) AddRouter(msgId uint16, router iface.Irouter) {
	s.MsgHandle.AddRouter(msgId,router)
}


func (s *Server) Start() {
	fmt.Printf("\033[45;36m%s\033[0m\n", "[cginx]"+utils.ServerOpt.Name)

	fmt.Println("[cginx]version:", utils.ServerOpt.Version)
	fmt.Println("[cginx]MaxPackageSize:", utils.ServerOpt.MaxPackageSize)
	fmt.Println("[Start] Listener at ip:", s.IP, "port:", s.Port)

	fmt.Println(s)
	if s.MsgHandle.Count() == 0 {
		panic("[msghandler is nil]")
	}

	go func() {

		addr, err := net.ResolveTCPAddr(s.IPType, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println(err)
			return
		}
		listen, err := net.ListenTCP(s.IPType, addr)

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("[Start] ", s.Name, " success, listening")
		var cid uint32
		cid = 0
		for {
			conn, err := listen.AcceptTCP()
			if err != nil {
				fmt.Println("Accept error:", err)
			}

			userConn := NewConnection(conn, cid, s.MsgHandle)
			cid++
			go userConn.Open()
		}
	}()

}

func (s *Server) Stop() {

}
func (s *Server) Serve() {

	//启动监听的go程
	s.Start()

	//启动之后可以加入其它的框架参数

	select {}

}
