package gateway

import (
	"net"
	"conf"
	"netlink"
	"io"
	"proto"
	"sync"
	"fmt"
)

type GateWay struct {
	cliser 		*netlink.Server
	serser 		*netlink.Server
	mc 			*proto.MessageCenter

	serEpList 	*Endpoint
}

func (gw *GateWay) registerMessage() {
	gw.mc = proto.NewMessageCenter()

	gw.mc.Register(100, &proto.TestData{})
}

func (gw *GateWay) Stop() {

}

/*
func (gw *GateWay) StartClientTest() {
	clientWait := new(sync.WaitGroup)
	clientWait.Add(1)
	for i := 0; i < 1; i++ {
		clientWait.Add(1)
		go func() {
			session, err := netlink.Dial("tcp",  "127.0.0.1:8899", netlink.ProtocolFunc(func (rw io.ReadWriter) (netlink.Codec, error) {
				return gw.ClientCodec(rw)
			}), 0)
			fmt.Println("create session errr", err)

			t := proto.TestData{A:10086}

			d := proto.Message{Cmd:100, Msg: t}
			err = session.Send(&d)
			fmt.Println("client send message ", d, err)


			m2, err := session.Receive()
			fmt.Println("Message Message ", m2, err)

			session.Close()
			clientWait.Done()
		}()
	}
	clientWait.Wait()
}
*/

func (gw *GateWay) Start() {

	gw.registerMessage()

	cfg := &conf.GatewayCfg{
		CliListenAddr: "127.0.0.1:8809",
		ClientSendChSize: 1024,

		SerListenAddr: "127.0.0.1:8890",
		ServerSendChSize: 1024,
	}

	var lns, lnc net.Listener
	var err error
	if lnc, err = net.Listen("tcp", cfg.CliListenAddr); err != nil {
		return
	}
	if lns, err = net.Listen("tcp", cfg.SerListenAddr); err != nil {
		return
	}

	go gw.StartClients(lnc, cfg)
	go gw.StartServers(lns, cfg)

	w := new(sync.WaitGroup)
	w.Add(1)
	w.Wait()
}

func (gw *GateWay) clientInit(io io.ReadWriter) (netlink.Codec, error) {
	return proto.NewCodec(gw.mc, 1, io.(net.Conn), 1024), nil
}

func (gw *GateWay) handleClientSession(session *netlink.Session) {
	defer session.Close()

	fmt.Println("server start")
	for {
		_, err := session.Receive()
		if err != nil {
			return
		}

		//err = session.Send(&d)
		if err != nil {
			return
		}
	}
}

func (gw *GateWay) StartClients(ln net.Listener, cfg *conf.GatewayCfg) {
	gw.cliser = netlink.NewServer(ln,
		netlink.ProtocolFunc(func (rw io.ReadWriter) (netlink.Codec, error) {
			return gw.clientInit(io)
		}),
		cfg.ClientSendChSize,
		netlink.HandlerFunc(func(session *netlink.Session) {
			gw.handleClientSession(session)
		}),
	)
	go gw.cliser.Serve()
}

func (gw *GateWay) serverInit(io io.ReadWriter) (netlink.Codec, error) {

	conn := io.(net.Conn)
	codec := proto.NewCodec(gw.mc, 1, conn, 0)

	msg, err := codec.Receive()
	if err != nil {
		conn.Close()
		return nil, nil
	}

	rsmsg, ok := msg.(*proto.MsgRegisterServer)
	if !ok {
		conn.Close()
		fmt.Println("assert type register server msg error")
		return nil, nil
	}

	ep := &ServerEndpoint{
		Id: rsmsg.Id,
		Type: rsmsg.Type,
	}
	if err := gw.serEpList.Add(ep); err != nil {
		codec.Send(&proto.MsgRegiserServerRes{Success: false})
		return nil, err
	}

	codec.Send(&proto.MsgRegiserServerRes{Success: true})

	return proto.NewCodec(gw.mc, uint32(ep.Id), conn, 1024), nil
}

func (gw *GateWay) handleServerSession(session *netlink.Session) {

}

func (gw *GateWay) StartServers(ln net.Listener, cfg *conf.GatewayCfg) {
	gw.serser = netlink.NewServer(
		ln,
		netlink.ProtocolFunc(func (rw io.ReadWriter) (netlink.Codec, error) {
			return gw.serverInit(rw)
		}),
		cfg.ServerSendChSize,
		netlink.HandlerFunc(func(session *netlink.Session) {
			gw.handleServerSession(session)
		}),
	)

	gw.serser.Serve()
}

func NewGateway() *GateWay {
	return &GateWay{
		serEpList: NewEndpoint(),
	}
}
