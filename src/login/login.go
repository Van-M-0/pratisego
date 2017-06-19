package login

type LoginServer struct {

}

func (ls *LoginServer) Start() {

}

func (ls *LoginServer) Stop() {

}


func NewLoginServer() *LoginServer {
	return &LoginServer{
	}
}