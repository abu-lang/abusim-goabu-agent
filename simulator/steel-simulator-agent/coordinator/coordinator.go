package coordinator

import (
	"net"
)

func GetConnection() (net.Conn, error) {

	tcpAddr, err := net.ResolveTCPAddr("tcp", "steel-coordinator:5001")
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
