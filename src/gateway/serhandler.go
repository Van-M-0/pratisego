package gateway

import (
	"netlink"
	"fmt"
	"proto"
	"errors"
)

type serverHandler struct {
	session 		*netlink.Session
	gw 				*GateWay
}

func newServerHandler(s *netlink.Session, g *GateWay) *serverHandler {
	return &serverHandler{
		session: s,
		gw: 	g,
	}
}

func (sh *serverHandler) init() error {

	data, err := sh.session.Receive()
	if err != nil {
		fmt.Println("gate way server hanlder init err", err)
		return err
	}

	msg, ok := data.(*proto.Message)
	if !ok {
		return errors.New("gate way server handler cast *Message")
	}

	registerMsg, ok := msg.Msg.(*proto.MsgRegisterServer)
	if !ok {
		return errors.New("gate way server handler cast *MsgRegisterServer")
	}

	if err := sh.gw.serEpList.Add(&ServerEndpoint{Id: registerMsg.Id, Type: registerMsg.Type}); err != nil {
		sh.session.Send(&proto.Message{Cmd: proto.Cmd_Register_Server_Res, Msg: &proto.MsgRegiserServerRes{
			Success: false,
		}})
		return err
	}

	sh.session.Send(&proto.Message{Cmd: proto.Cmd_Register_Server_Res, Msg: &proto.MsgRegiserServerRes{
		Success: true,
	}})

	return nil
}

func (sh *serverHandler) close() {
	sh.session.Close()
}

func (sh *serverHandler) start() {
	for {
		data, err := sh.session.Receive()
		if err != nil {
			sh.close()
			return
		}

		msg, ok := data.(*proto.Message)
		if !ok {
			fmt.Println("gate way server handler cast *Message err")
			continue
		}

		if err := sh.handle(msg.Cmd, msg.Msg); err != nil {
			sh.close()
			return
		}
	}
}

func (sh *serverHandler) handle(cmd uint32, msg interface{}) error {
	return nil
}

