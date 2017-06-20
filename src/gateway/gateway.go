package gateway

import (
	"net"
	"conf"
	"netlink"
	"io"
	"proto"
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

	proto.RegisterGatewayMsg(gw.mc)
}

func (gw *GateWay) Stop() {

}

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
}

func (gw *GateWay) clientInit(io io.ReadWriter) (netlink.Codec, error) {
	return proto.NewCodec(gw.mc, 1, io.(net.Conn), 1024), nil
}

func (gw *GateWay) handleClientSession(session *netlink.Session) {
	ch := newClientHandler(session, gw)
	if ch.init() != nil {
		fmt.Println("init client handler err")
	}
	ch.start()
}

func (gw *GateWay) StartClients(ln net.Listener, cfg *conf.GatewayCfg) {
	gw.cliser = netlink.NewServer(ln,
		netlink.ProtocolFunc(func (rw io.ReadWriter) (netlink.Codec, error) {
			return gw.clientInit(rw)
		}),
		cfg.ClientSendChSize,
		netlink.HandlerFunc(func(session *netlink.Session) {
			gw.handleClientSession(session)
		}),
	)
	go gw.cliser.Serve()
}

func (gw *GateWay) handleServerSession(session *netlink.Session) {
	sh := newServerHandler(session, gw)
	if sh.init() != nil {
		return
	}
	sh.start()
}

func (gw *GateWay) StartServers(ln net.Listener, cfg *conf.GatewayCfg) {
	gw.serser = netlink.NewServer(
		ln,
		netlink.ProtocolFunc(func (rw io.ReadWriter) (netlink.Codec, error) {
			return proto.NewCodec(gw.mc, uint32(0), rw.(net.Conn), 1024), nil
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
