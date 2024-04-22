package handlers

import (
	"net"

	"github.com/sagini18/message-broker/consumer/internal/types"
)

func InitConnection(tcpConn *types.TcpConn) error {
	conn, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		return err
	}
	tcpConn.New(conn)
	return nil
}
