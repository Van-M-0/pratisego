package main

/*
type message struct {
	a 		uint16
}

func newMsgCodec(rw io.ReadWriter) (netlink.Codec, error) {
	dec := &proto.MessageDec{}
	dec.Io(rw.(io.ReadWriteCloser))
	return dec, nil
}

func messagetest(session *netlink.Session) {
	main, sub := 10, 100
	m1 := &message{a:101}
	session.Send(&proto.Message{uint8(main), uint8(sub), m1})

	m2, err := session.Receive()
	fmt.Println("message test ", m2, err)
}

func MessageSync() {
	main, sub := 10, 100
	type message struct {
		a 		uint16
	}
	m1 := &message{a:101}
	proto.RegisterMessage(uint8(main), uint8(sub), m1)

	session_test(0, messagetest)
}

func MessageAsync() {
	session_test(1024, messagetest)
}
*/
