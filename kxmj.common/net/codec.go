package net

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/panjf2000/gnet/v2"
	"google.golang.org/protobuf/proto"
)

const (
	MagicNumber       = uint16(13145)
	InnerHeaderLength = 16
	OuterHeaderLength = 14
)

func bytesOrder() binary.ByteOrder {
	return binary.BigEndian
}

// 包头定义(内网)
//
// * 0           2                       6           8           10          12                      16
// * +-----------+-----------------------+-----------+-----------+-----------+-----------------------+
// * |   magic   |       body len        |  msg id   |  svr_type |   svr_id  |       user_id         |
// * +-----------+-----------+-----------+-----------+-----------+-----------+-----------+-----------+
// * |                                                                                               |
// * +                                                                                               +
// * |                                            body bytes                                         |
// * +                                                                                               +
// * |                                             ... ...                                           |
// * +-----------------------------------------------------------+-----------+-----------------------+

func (msg *Message) Encode() ([]byte, error) {
	length := uint32(len(msg.Data))
	buffer := bytes.NewBuffer(make([]byte, 0, InnerHeaderLength+length))
	if err := binary.Write(buffer, bytesOrder(), MagicNumber); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, bytesOrder(), length); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, bytesOrder(), msg.MsgId); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, bytesOrder(), msg.SvrType); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, bytesOrder(), msg.SvrId); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, bytesOrder(), msg.UserId); err != nil {
		return nil, err
	}

	if length > 0 {
		if err := binary.Write(buffer, bytesOrder(), msg.Data); err != nil {
			return nil, err
		}
	}

	return buffer.Bytes(), nil
}

func (msg *Message) Decode(body proto.Message) error {
	return proto.Unmarshal(msg.Data, body)
}

func Marshal(body proto.Message) []byte {
	data, err := proto.Marshal(body)
	if err != nil {
		return nil
	}
	return data
}

func Unpack(data []byte) (*Message, error) {
	if InnerHeaderLength > len(data) {
		return nil, fmt.Errorf("invalid msg List")
	}

	magicNumber := bytesOrder().Uint16(data[0:2])
	if magicNumber != MagicNumber {
		return nil, fmt.Errorf("invalid MagicNumber")
	}

	length := bytesOrder().Uint32(data[2:6])
	if len(data) != int(length)+InnerHeaderLength {
		return nil, fmt.Errorf("is invalid List")
	}

	msg := &Message{
		MsgId:   bytesOrder().Uint16(data[6:8]),
		SvrType: bytesOrder().Uint16(data[8:10]),
		SvrId:   bytesOrder().Uint16(data[10:12]),
		UserId:  bytesOrder().Uint32(data[12:16]),
	}

	if len(data) > InnerHeaderLength {
		msg.Data = data[InnerHeaderLength:]
	}

	return msg, nil
}

// 包头定义(外网)
//
// * 0           2                       6           8          10          12          14
// * +-----------+-----------------------+-----------+-----------+-----------+-----------+
// * |   magic   |       body len        |  msg id   |  svr_type |   svr_id  |     ex    |
// * +-----------+-----------+-----------+-----------+-----------+-----------+-----------+
// * |                                                                                   |
// * +                                                                                   +
// * |                                  body bytes                                       |
// * +                                                                                   +
// * |                                   ... ...                                         |
// * +-----------------------------------------------------------+-----------+-----------+

var ErrIncompletePacket = errors.New("incomplete packet")
var ErrMagicNumberPacket = errors.New("invalid MagicNumber")

type OuterCodec struct {
}

func (oc *OuterCodec) Encode(msg *Message) ([]byte, error) {
	length := uint32(len(msg.Data))
	buffer := bytes.NewBuffer(make([]byte, 0, OuterHeaderLength+length))
	if err := binary.Write(buffer, bytesOrder(), MagicNumber); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, bytesOrder(), length); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, bytesOrder(), msg.MsgId); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, bytesOrder(), msg.SvrType); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, bytesOrder(), msg.SvrId); err != nil {
		return nil, err
	}
	ex := uint16(0)
	if err := binary.Write(buffer, bytesOrder(), ex); err != nil {
		return nil, err
	}

	if length > 0 {
		if err := binary.Write(buffer, bytesOrder(), msg.Data); err != nil {
			return nil, err
		}
	}

	return buffer.Bytes(), nil
}

func (oc *OuterCodec) Decode(conn gnet.Conn) (*Message, error) {
	data, _ := conn.Peek(OuterHeaderLength)
	if len(data) < OuterHeaderLength {
		return nil, ErrIncompletePacket
	}

	magicNumber := bytesOrder().Uint16(data[0:2])
	if magicNumber != MagicNumber {
		return nil, ErrMagicNumberPacket
	}

	length := bytesOrder().Uint32(data[2:6])
	msgLen := OuterHeaderLength + int(length)
	if conn.InboundBuffered() < msgLen {
		return nil, ErrIncompletePacket
	}

	data, _ = conn.Peek(msgLen)
	_, _ = conn.Discard(msgLen)

	if len(data) < int(length)+OuterHeaderLength {
		return nil, ErrIncompletePacket
	}

	msg := &Message{
		MsgId:   bytesOrder().Uint16(data[6:8]),
		SvrType: bytesOrder().Uint16(data[8:10]),
		SvrId:   bytesOrder().Uint16(data[10:12]),
	}

	if len(data) > OuterHeaderLength {
		msg.Data = data[OuterHeaderLength:msgLen]
	}

	return msg, nil
}
