package gateway

import (
	"netlink"
	"fmt"
	"proto"
)

type clientHandler struct {
	session 		*netlink.Session
	gw 				*GateWay
}

func newClientHandler(session *netlink.Session, gw *GateWay) *clientHandler {
	return &clientHandler{
		session: session,
		gw: gw,
	}
}

func (ch *clientHandler) init() error {
	return nil
}

func (ch *clientHandler) close() {
	ch.session.Close()
}

func (ch *clientHandler) start() {
	for {
		data, err := ch.session.Receive()
		if err != nil {
			ch.close()
			return
		}

		msg, ok := data.(*proto.Message)
		if !ok {
			fmt.Println("client handler cast *Message err", msg)
			continue
		}

		if err := ch.handle(msg.Cmd, msg.Msg); err != nil {
			ch.close()
			return
		}
	}
}

func (ch *clientHandler) handle(cmd uint32, msg interface{}) error {
	return nil
}

