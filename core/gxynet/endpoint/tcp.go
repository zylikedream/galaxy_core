package endpoint

import (
	"bytes"
	"context"
	"net"

	"github.com/zylikedream/galaxy/core/gxynet/message"
	"github.com/zylikedream/galaxy/core/gxynet/processor"
)

type TcpEndpoint struct {
	conn net.Conn
	proc processor.Processor
	buf  *bytes.Buffer
	data any
}

func NewTcpEndPoint(conn net.Conn, proc processor.Processor) *TcpEndpoint {
	return &TcpEndpoint{
		proc: proc,
		conn: conn,
		buf:  &bytes.Buffer{},
	}
}

func (t *TcpEndpoint) DecodeMsg(data []byte) (*message.Message, error) {
	t.buf.Write(data)
	pkgLen, msg, err := t.proc.Decode(t.buf.Bytes())
	if err != nil {
		return nil, err
	}
	if msg == nil {
		return nil, nil
	}
	t.buf.Next(int(pkgLen))
	return msg, nil
}

func (t *TcpEndpoint) Send(msg *message.Message) error {
	data, err := t.proc.Encode(msg)
	if err != nil {
		return err
	}
	if _, err = t.conn.Write(data); err != nil {
		return err
	}
	return nil
}

func (t *TcpEndpoint) Close(ctx context.Context, err error) {
	t.conn.Close()
}

func (t *TcpEndpoint) Conn() net.Conn {
	return t.conn
}

func (t *TcpEndpoint) GetData() any {
	return t.data
}

func (t *TcpEndpoint) SetData(d any) {
	t.data = d
}
