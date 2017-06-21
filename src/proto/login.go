package proto

func RegisterGatewayMsg(mc *MessageCenter) {
	RegisterCommonMsg(mc)
}

const CMD_LOGINLOGIN = CMD_login + 1
type MsgLoginLogin struct {
	Name 			string
	UserId 			uint32
}

const CMD_LOGIN_LOGIN_RES = CMD_login + 2
type MsgLoginLoginRes struct {
	ErrCode 		int
}

