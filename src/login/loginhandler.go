package login

import (
	"netlink"
	"net"
	"proto"
	"errors"
	"fmt"
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
		mc: proto.NewMessageCenter(),
	}
}

func (lh *loginHandler) init() error {

	proto.RegisterLoginMsg(lh.mc)

	if err := lh.connectGateway(); err != nil {
		lh.close()
		return err
	}

	return nil
}

func (lh *loginHandler) close() {
	fmt.Println("login handler close")
	lh.session.Close()
}

func (lh *loginHandler) start() error {
	for {
		imsg, err := lh.session.Receive()
		if err != nil {
			return err
		}

		msg, ok := imsg.(*proto.Message)
		if !ok {
			return errors.New("login handler cast *Message err")
		}

		if err := lh.handleMessage(int(msg.Cmd), &msg.Msg); err != nil {
			return err
		}
	}
}

func (lh *loginHandler) connectGateway() error {
	conn, err := net.Dial("tcp", lh.GatewayAddr)
	if err != nil {
		return err
	}

	codec := proto.NewCodec(lh.mc, uint32(lh.ServerId), conn.(net.Conn), lh.RecvBuffSize)
	lh.session = netlink.NewSession(codec, lh.SendChSize)

	if err := lh.register2remote(); err != nil {
		return err
	}
	return nil
}

func (lh *loginHandler) register2remote() error {

	if err := lh.SendMsg(proto.Cmd_Register_Server, &proto.MsgRegisterServer{Type: 1, Id: lh.ServerId}); err != nil {
		return err
	}

	data, err := lh.session.Receive()
	if err != nil {
		return err
	}

	imsg, ok := data.(*proto.Message)
	if !ok {
		return errors.New("login handler cast *Message err")
	}

	msg, ok := imsg.Msg.(*proto.MsgRegiserServerRes)
	if !ok {
		return errors.New("login handler cast *MsgregiserServerRes err")
	}

	fmt.Println("register res ", msg)
	if !msg.Success {
		return errors.New("cannot register login handler to gateway")
	}

	return nil
}


func (lh *loginHandler) SendMsg(cmd int, data interface{}) error {
	return lh.session.Send(&proto.Message{Cmd: uint32(cmd), Msg: data})
}

func (lh *loginHandler) handleMessage(cmd int, data *interface{}) error {
	return nil
}
