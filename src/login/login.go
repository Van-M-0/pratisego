package login

import "fmt"

type LoginServer struct {
	lghandler 			*loginHandler
}

func (ls *LoginServer) Start() error {

	defer ls.Stop()

	ls.lghandler.ServerId  = 1
	ls.lghandler.GatewayAddr = "127.0.0.1:8890"

	if err := ls.lghandler.init(); err != nil {
		fmt.Println("init login handler err ", err)
		return err
	}

	if err := ls.lghandler.start(); err != nil {
		fmt.Println("start login handler err ", err)
		return err
	}

	return nil
}

func (ls *LoginServer) Stop() {

}

func NewLoginServer() *LoginServer {
	return &LoginServer{
		lghandler: newLoginHandler(),
	}
}