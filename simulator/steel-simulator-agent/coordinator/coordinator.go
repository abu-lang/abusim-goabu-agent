package coordinator

import (
	"errors"
	"net"
	"steel-simulator-config/communication"
)

type Coordinator struct {
	conn  *net.TCPConn
	coord *communication.Coordinator
}

func New() (*Coordinator, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "steel-coordinator:5001")
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}

	return &Coordinator{
		conn:  conn,
		coord: communication.New(conn),
	}, nil
}

func (c *Coordinator) SendInit(configStr string) error {
	err := c.coord.Write(&communication.CoordinatorMessage{
		Type:    communication.CoordinatorMessageTypeINIT,
		Payload: configStr,
	})
	if err != nil {
		return err
	}
	msg, err := c.coord.Read()
	if err != nil {
		return err
	}
	if msg.Type != communication.CoordinatorMessageTypeACK {
		return errors.New("unexpected response to init")
	}
	return nil
}

func (c *Coordinator) Close() {
	c.conn.Close()
}
