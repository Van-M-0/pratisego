package proto

const (
	CMD_common 			= 0
	CMD_login 			= 1000
	CMD_gateway 		= 2000
	CMD_gameserver 		= 3000
	CMD_master 			= 5000
)

type Message struct {
	Cmd 		uint32
	Msg 		interface{}
}

type TestData struct {
	A 			int
}

const Cmd_Register_Server = CMD_common + 1
type MsgRegisterServer struct {
	Type 		int
	Id 			int
}

const Cmd_Register_Server_Res = CMD_common + 2
type MsgRegiserServerRes struct {
	Success 	bool
}

func RegisterCommonMsg(mc *MessageCenter) {
	mc.Register(Cmd_Register_Server, &MsgRegisterServer{})
	mc.Register(Cmd_Register_Server_Res, &MsgRegiserServerRes{})
}
