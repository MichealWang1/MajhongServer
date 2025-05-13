package net

import (
	"bufio"
	"errors"
	"net"
)

type TcpClient struct {
	addr     string
	conn     net.Conn
	callback func(msg *Message)
	msg      []byte
	oc       *OuterCodec
}

func NewTcpClient(addr string, callback func(msg *Message)) *TcpClient {
	client := &TcpClient{
		addr:     addr,
		callback: callback,
		msg:      make([]byte, 0),
		oc:       &OuterCodec{},
	}
	return client
}

func (tc *TcpClient) Connect() error {
	conn, err := net.Dial("tcp", tc.addr)
	if err != nil {
		return err
	}
	tc.conn = conn
	go tc.read()
	return nil
}

func (tc *TcpClient) read() {
	reader := bufio.NewReader(tc.conn)
	for {
		buf := make([]byte, 1024)
		n, err := reader.Read(buf)
		if err != nil {
			return
		}

		if n == 0 {
			return
		}

		tc.msg = append(tc.msg, buf[:n]...)
		if OuterHeaderLength > len(tc.msg) {
			continue
		}

		msgList := make([]*Message, 0)
		startIndex := 0
		for {
			if startIndex >= len(tc.msg)-1 {
				break
			}

			length := bytesOrder().Uint32(tc.msg[startIndex+2 : startIndex+6])
			if int(length) > len(tc.msg)-startIndex {
				break
			}

			msg := &Message{
				MsgId:   bytesOrder().Uint16(tc.msg[startIndex+6 : startIndex+8]),
				SvrType: bytesOrder().Uint16(tc.msg[startIndex+8 : startIndex+10]),
				SvrId:   bytesOrder().Uint16(tc.msg[startIndex+10 : startIndex+12]),
			}

			if length > 0 {
				msg.Data = tc.msg[startIndex+OuterHeaderLength : startIndex+int(OuterHeaderLength+length)]
			}

			msgList = append(msgList, msg)
			startIndex += int(OuterHeaderLength + length)
		}

		tc.msg = tc.msg[startIndex:]

		if tc.callback != nil {
			for _, msg := range msgList {
				tc.callback(msg)
			}
		}
	}
}

func (tc *TcpClient) Send(msg *Message) error {
	data, err := tc.oc.Encode(msg)
	if err != nil {
		return err
	}

	n, err := tc.conn.Write(data)
	if err != nil {

		return err
	}

	if n <= 0 {
		return errors.New("write zero data, will reconnect")
	}

	return nil
}

func (tc *TcpClient) Close() error {
	return tc.conn.Close()
}
