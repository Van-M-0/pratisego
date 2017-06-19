package main

import (
	"io"
	"netlink"
	"encoding/binary"
	"fmt"
	"time"
	"bytes"
	"math/rand"
)

func NewTestCodec(rw io.ReadWriter) (netlink.Codec, error) {
	return &TestCodec{
		rw: rw.(io.ReadWriteCloser),
	}, nil
}

type TestCodec struct {
	rw io.ReadWriteCloser
}

func (c *TestCodec) Send(msg interface{}) error {
	var head [2]byte
	binary.LittleEndian.PutUint16(head[:], uint16(len(msg.([]byte))))
	_, err := c.rw.Write(head[:])
	if err != nil {
		return err
	}
	_, err = c.rw.Write(msg.([]byte))
	if err != nil {
		return err
	}
	return nil
}

func (c *TestCodec) Receive() (interface{}, error) {
	var head [2]byte
	_, err := io.ReadFull(c.rw, head[:])
	if err != nil {
		return nil, err
	}
	n := binary.LittleEndian.Uint16(head[:])
	buf := make([]byte, n)
	_, err = io.ReadFull(c.rw, buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (c *TestCodec) Close() error {
	return c.rw.Close()
}

func (c *TestCodec) ClearSendChan(ch <-chan interface{}) {
	for _ = range ch {
	}
}

func RandBytes(n int) []byte {
	n = rand.Intn(n) + 1
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte(rand.Intn(255))
	}
	return b
}

func BytesTest(session *netlink.Session) {
	for i := 0; i < 2; i++ {
		fmt.Println("sleep ...")
		time.Sleep(5*time.Second)
		msg1 := RandBytes(512)
		err := session.Send(msg1)
		if err != nil {
			fmt.Println("byte test errr")
		}

		msg2, err := session.Receive()
		bytes.Equal(msg1, msg2.([]byte))
	}
}

func Test_Sync() {
	//session_test(0, BytesTest)
}

func Test_Async() {
	//session_test(1024, BytesTest)
}
