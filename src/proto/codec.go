package proto

import (
	"net"
	"bufio"
	"io"
	"encoding/binary"
	"errors"
	"gopkg.in/vmihailenco/msgpack.v2"
)

const SizeOfHeadDataLen = 2
const SizeOfCmdDataLen = 4
const SizeOfHeaderLen  = SizeOfHeadDataLen + SizeOfCmdDataLen

const (
	MaxPacketSize		= 0xffff
)

var (
	ErrOverMaxPacketSize = errors.New("over max packet size")
	ErrConvertSendMsg	 = errors.New("convert msg error")
)

type Codec struct {
	id 			uint32
	conn 		net.Conn
	reader 		*bufio.Reader
	headBuf 	[]byte
	headData 	[SizeOfHeaderLen]byte
	sendBuf 	[MaxPacketSize]byte
	mc 			*MessageCenter
}

func NewCodec(mc *MessageCenter, id uint32, conn net.Conn, bufferSize int) *Codec {
	c := &Codec{
		id: 		id,
		conn: 		conn,
		reader:		bufio.NewReaderSize(conn, bufferSize),
		mc: 		mc,
	}
	c.headBuf = c.headData[:]
	return c
}

func (c *Codec) Receive() (interface{}, error) {
	if _, err := io.ReadFull(c.conn, c.headBuf); err != nil {
		return nil, err
	}

	length := int(binary.LittleEndian.Uint16(c.headBuf[:SizeOfHeadDataLen]))
	cmd := uint32(binary.LittleEndian.Uint32(c.headBuf[SizeOfHeadDataLen:]))
	if length > MaxPacketSize {
		return nil, ErrOverMaxPacketSize
	}

	buf := make([]byte, length)
	if _, err := io.ReadFull(c.reader, buf); err != nil {
		return nil, err
	}

	msg, err := c.mc.NewMessage(cmd)
	if err != nil {
		return nil, err
	}
	err = msgpack.Unmarshal(buf, &msg)
	return &msg, err
}

func (c *Codec) Send(i interface{}) error {

	data, ok := i.(*Message)
	if !ok {
		return ErrConvertSendMsg
	}

	_, err := c.mc.GetType(data.Cmd)
	if err != nil {
		return err
	}

	buf, err := msgpack.Marshal(data.Msg)
	if err != nil {
		return err
	}

	binary.LittleEndian.PutUint16(c.sendBuf[:], uint16(len(buf)))
	binary.LittleEndian.PutUint32(c.sendBuf[SizeOfHeadDataLen:], data.Cmd)

	copy(c.sendBuf[SizeOfHeaderLen:], buf)

	totalLen := SizeOfHeaderLen + len(buf)
	_, err = c.conn.Write(c.sendBuf[:totalLen])
	return err
}

func (c *Codec) Close() error {
	return c.conn.Close()
}
