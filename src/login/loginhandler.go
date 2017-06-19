package login

import (
	"netlink"
	"net"
	"proto"
	"errors"
)

type loginHandler struct {
	ServerId 			int
	GatewayAddr 		string
	SendChSize 			int
	RecvBuffSize		int
	session 			*netlink.Session
	mc 					*proto.MessageCenter
}

func newLoginHandler() *loginHandler {
	return &loginHandler{

	}
}

func (lh *loginHandler) registerMessage () {

}

func (lh *loginHandler) register2remote(conn net.Conn) error {
	codec := proto.NewCodec(lh.mc, 1, conn, 0)
	codec.Send(&proto.MsgRegisterServer{Type:1, Id:lh.ServerId})

	msg, err := codec.Receive()
	if err != nil {
		return err
	}

	rsmsg, ok := msg.(*proto.MsgRegiserServerRes)
	if !ok {
		return errors.New("assert type register server msg error")
	}

	if !rsmsg.Success {
		return errors.New("cannot register login handler to gateway")
	}

	return nil
}

func (lh *loginHandler) connectGateway(addr string) error {

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	codec := proto.NewCodec(lh.mc, 1, conn.(net.Conn), lh.RecvBuffSize)
	lh.session = netlink.NewSession(codec, lh.SendChSize)

	if err := lh.register2remote(conn); err != nil {
		conn.Close()
		return err
	}

	return nil
}

func (lh *loginHandler) Start(addr string) error {
	lh.registerMessage()

	if err := lh.connectGateway(addr); err != nil {
		return err
	}

	go lh.handle(lh.session)

	return nil
}

func (lh *loginHandler) Send(cmd int, data interface{}) {
	lh.session.Send(&proto.Message{Cmd: uint32(cmd), Msg: data})
}

func (lh *loginHandler) handle(session *netlink.Session) {
	defer func() {
		lh.session.Close()
	}()

	for {
		imsg, err := session.Receive()
		if err != nil {
			return
		}

		msg, ok := imsg.(*proto.Message)
		if !ok {
			return
		}

		lh.handleMessage(int(msg.Cmd), &msg.Msg)
	}
}

func (lh *loginHandler) handleMessage(cmd int, data *interface{}) {

}