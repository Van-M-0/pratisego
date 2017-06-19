package main

import (
	"gateway"
)

/*
import (
	"netlink"
	"fmt"
	"sync"
	"proto"
	"net"
)

func session_test(sendsize int, test func(session *netlink.Session), codec func(conn net.Conn) (netlink.Codec, error)) {
	server, _ := netlink.Listen("tcp", "127.0.0.1:47474", netlink.ProtocolFunc(codec), sendsize, netlink.HandlerFunc(func (session *netlink.Session) {
		defer session.Close()

		fmt.Println("server start")
		for {
			msg, err := session.Receive()
			fmt.Println("server receive msg", msg)
			if err != nil {
				fmt.Println("session receive msg ", err)
				return
			}
			fmt.Println("server send msg", msg)
			err = session.Send(msg)
			if err != nil {
				fmt.Println("session send msg ", msg)
				return
			}
		}
	}))

	go server.Serve()

	addr := server.Listener().Addr().String()
	fmt.Println("listen addres is ", addr)

	clientWait := new(sync.WaitGroup)
	for i := 0; i < 1; i++ {
		clientWait.Add(1)
		go func() {
			session, err := netlink.Dial("tcp", addr, netlink.ProtocolFunc(codec), sendsize)
			fmt.Println("create session errr", err)
			test(session)
			session.Close()
			clientWait.Done()
		}()
	}
	clientWait.Wait()

	server.Stop()
}

func main() {
	//Test_Sync()
	//MessageSync()

	struct test {

	}

	mc := proto.NewMessageCenter()
	mc.Register(100)


	//gw := gateway.NewGateway()
	//gw.Start()
}
*/


func main() {

	msgpack_test()

	gw := gateway.NewGateway()
	gw.Start()
}
