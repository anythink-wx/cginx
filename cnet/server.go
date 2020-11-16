package cnet

import (
	"cginx/iface"
	"fmt"
	"net"
)

type Server struct {
	Name   string
	IPType string
	IP     string
	Port   int
}

func NewServer(name string, ) iface.Iserver{
	return &Server{
		Name:   name,
		IPType: "tcp",
		IP:     "0.0.0.0",
		Port:   8810,
	}
}

func callback(conn *net.TCPConn, data []byte, ln int) (err error) {
	_, err = conn.Write(data[:ln])
	if err != nil {
		fmt.Println("write error:", err)
	}
	return
}

func (s *Server) Start() {
	fmt.Println("[Start] hello:", s.Name)
	fmt.Println("[Start] Listener at ip:", s.IP, "port:", s.Port)

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

			userConn := NewConnection(conn, cid, callback)
			cid++
			go userConn.open()
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