package types

import (
	"net"
)

type TcpConn struct {
	Conn net.Conn
}

func (t *TcpConn) New(conn net.Conn) {
	t.Conn = conn
}
