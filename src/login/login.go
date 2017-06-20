package login

type LoginServer struct {
	lghandler 			*loginHandler
}

func (ls *LoginServer) Start() error {

	ls.lghandler.ServerId  = 1
	ls.lghandler.GatewayAddr = "127.0.0.1:8890"

	if err := ls.lghandler.Start(ls.lghandler.GatewayAddr); err != nil {
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