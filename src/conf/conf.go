package conf


type GatewayCfg struct {
	CliListenAddr				string
	ClientSendChSize 			int

	SerListenAddr 				string
	ServerSendChSize			int
}
