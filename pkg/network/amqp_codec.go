package network

import (
	"amqp_go/pkg/codec"
	"encoding/binary"
	"errors"
	"github.com/panjf2000/gnet"
)

type amqpCodec struct {
}

func (a *amqpCodec) Encode(c gnet.Conn, buf []byte) ([]byte, error) {
	if codec.IsProtocolHeader(buf) {
		return buf, nil
	}
	length := len(buf) + 4
	out := make([]byte, 4)
	binary.BigEndian.PutUint32(out, uint32(length))
	return append(out, buf...), nil
}

type innerBuffer []byte

func (in *innerBuffer) readN(n int) (buf []byte, err error) {
	if n == 0 {
		return nil, nil
	}

	if n < 0 {
		return nil, errors.New("negative length is invalid")
	} else if n > len(*in) {
		return nil, errors.New("exceeding buffer length")
	}
	buf = (*in)[:n]
	*in = (*in)[n:]
	return
}

func (a amqpCodec) Decode(c gnet.Conn) ([]byte, error) {
	var (
		in     innerBuffer
		header []byte
		err    error
	)
	in = c.Read()
	if len(in) < 8 {
		return nil, errors.New("incomplete packet")
	}
	if codec.IsProtocolHeader(in[0:8]) {
		c.ShiftN(8)
		return in[0:8], nil
	}

	lenBuf, frameLength, err := a.getUnadjustedFrameLength(&in)
	if err != nil {
		return nil, err
	}

	// real message length
	msgLength := int(frameLength) - 4
	msg, err := in.readN(msgLength)
	if err != nil {
		return nil, errors.New("incomplete packet")
	}

	fullMessage := make([]byte, len(header)+len(lenBuf)+msgLength)
	copy(fullMessage, header)
	copy(fullMessage[len(header):], lenBuf)
	copy(fullMessage[len(header)+len(lenBuf):], msg)
	c.ShiftN(len(fullMessage))
	return fullMessage[4:], nil
}

func (a *amqpCodec) getUnadjustedFrameLength(in *innerBuffer) ([]byte, uint64, error) {
	lenBuf, err := in.readN(4)
	if err != nil {
		return nil, 0, errors.New("unexpected eof")
	}
	return lenBuf, uint64(binary.BigEndian.Uint32(lenBuf)), nil
}
