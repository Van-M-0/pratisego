package proto

/*
import (
	"reflect"
	"io"
	"gopkg.in/vmihailenco/msgpack.v2"
	"fmt"
	"encoding/binary"
	"net"
)

var (
	types 	map[uint32]reflect.Type
	ids 	map[reflect.Type]uint32
)

func init() {
	types = make(map[uint32]reflect.Type)
	ids = make(map[reflect.Type]uint32)
}

func RegisterMessage(main, sub uint8, i interface{}) {
	id := uint32(CombineMessageId(main, sub))
	rt := reflect.TypeOf(i)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	types[id] = rt
	ids[rt] = id
}

func NewRegisterMessage(id uint32) (interface{}, error) {
	if t, ok := types[id]; ok {
		return reflect.New(t).Interface(), nil
	} else {
		err := fmt.Errorf("not found %d", id)
		return nil, err
	}
}

func ParseMessageId(id uint16) (uint8, uint8) {
	return uint8(id >> 8), uint8(id & 0xffff)
}

func CombineMessageId(main, sub uint8) uint16 {
	a := uint16(main) << 8
	b := a | uint16(sub)
	return uint16(b)
}

type MessageDec	struct {
	rwc 			io.ReadWriteCloser
}

func (dec *MessageDec) Io(o io.ReadWriteCloser) {
	dec.rwc = o
}

func (dec *MessageDec) GetIo() net.Conn {
	conn, ok := dec.rwc.(net.Conn)
	if !ok {
		return nil
	}
	return conn
}

func (dec *MessageDec) allocHeader(len uint16, main, sub uint8) []byte {
	b := make([]byte, MESSAGEHEADER_LEN)
	binary.LittleEndian.PutUint16(b, len)
	b[2] = main
	b[3] = sub
	return b
}

func (dec *MessageDec) toHeader(b []byte) (uint16, uint8, uint8) {
	return binary.LittleEndian.Uint16(b), b[2], b[3]
}

func (dec *MessageDec) Receive() (interface{}, error) {
	fmt.Println("codec receive ")

	var header [MESSAGEHEADER_LEN]byte
	_, err := io.ReadFull(dec.rwc, header[:])
	if err != nil {
		return nil, err
	}

	len, main, sub := dec.toHeader(header[:])
	id := CombineMessageId(main, sub)

	body := make([]byte, len)
	_, err = io.ReadFull(dec.rwc, body)
	if err != nil {
		return nil, err
	}

	fmt.Println("receieve msg header body", header, body)

	fmt.Println("id is : ", id, types[uint32(id)], reflect.ValueOf(types[uint32(id)]))

	imsg, err := NewRegisterMessage(uint32(id))
	fmt.Println("new rigsiter message ", imsg, err)
	if err != nil {
		return nil, err
	}

	err = msgpack.Unmarshal(body, &imsg)
	if err != nil {
		return nil, err
	}
	fmt.Println("unmarshal msg ", imsg)

	return &Message{main, sub, imsg}, nil
}

func (dec *MessageDec) Send(msg interface{}) error {
	mp, ok := msg.(*Message)
	if !ok {
		fmt.Println("type assertion failed ", ok)
	}

	body, err := msgpack.Marshal(mp.Msg)
	if err != nil {
		return err
	}

	headBuf := dec.allocHeader(uint16(len(body)), mp.Main, mp.Sub)
	_, err = dec.rwc.Write(headBuf)
	if err != nil {
		return err
	}

	_, err = dec.rwc.Write(body)
	if err != nil {
		return err
	}

	fmt.Println("codec send header, body", headBuf, body)
	return nil
}

func (dec *MessageDec) Close() error {
	dec.rwc.Close()
	return nil
}
*/