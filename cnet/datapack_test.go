package cnet

import (
	"fmt"
	"net"
	"testing"
)

func TestDataPack_Pack(t *testing.T) {

	listen, err2 := net.Listen("tcp", "127.0.0.1:7777")
	if err2 != nil {
		panic(err2)
	}

	go func() {

		for {
			conn, err := listen.Accept()
			if err != nil {
				fmt.Println(err)
				return
			}

			go func(conn net.Conn) {
				pack := NewDataPack()

				for {

					headData := make([]byte, pack.GetHeadLen())
					ln, err := conn.Read(headData)
					if err != nil {
						fmt.Println(err)
						break
					}
					fmt.Println("read ln", ln)
					msgHead, err := pack.Unpack(headData[:ln])
					if err != nil {
						fmt.Println("server unpack head err=", err)
						return
					}
					if msgHead.GetMsgLen() > 0 {
						//说明msg有数据
						msg := msgHead.(*message)
						msg.Data = make([]byte, msgHead.GetMsgLen())
						_, err := conn.Read(msg.Data)
						if err != nil {
							fmt.Println("server unpack data err=", err)
							return
						}

						fmt.Println("recive msg:", msg.GetMsgId(), "data len:", msg.GetMsgLen(), "data:", string(msg.GetMsgData()))
					}

				}
				//msg, err := pack.Unpack(buf)
				//if err != nil {
				//	fmt.Println(err)
				//	return
				//}
				//msg.GetMsgLen()

			}(conn)

		}

	}()

	conn, err2 := net.Dial("tcp", "127.0.0.1:7777")
	if err2 != nil {
		panic("client dial err"  + err2.Error())
	}
	dataPack := NewDataPack()
	//模拟粘包过程

	msg1 := &message{
		Id: 1,
		DataLen: 5,
		Data: []byte("hello"),
	}
	sendPkg1, err2 := dataPack.Pack(msg1)
	if err2 != nil {
		fmt.Print("pack error", sendPkg1)
	}


	msg2 := &message{
		Id: 2,
		DataLen: 11,
		Data: []byte("hello world"),
	}
	sendPkg2, err2 := dataPack.Pack(msg2)
	if err2 != nil {
		fmt.Print("pack error", sendPkg1)
	}

	send := append(sendPkg1,sendPkg2...)
	conn.Write(send)

	fmt.Println(send)

	select {}





}
