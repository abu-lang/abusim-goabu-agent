package coordinator

import (
	"errors"
	"log"
	"net"
	"steel-lang/semantics"
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

func (c *Coordinator) SendSelfName(name string) error {
	err := c.coord.Write(&communication.CoordinatorMessage{
		Type:    communication.CoordinatorMessageTypeINIT,
		Payload: name,
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

func (c *Coordinator) HandleMessages(exec *semantics.MuSteelExecuter) {
	for {
		msg, err := c.coord.Read()
		if err != nil {
			log.Println(err)
			break
		}
		switch msg.Type {
		case communication.CoordinatorMessageTypeMemoryREQ:
			err := c.coord.Write(&communication.CoordinatorMessage{
				Type:    communication.CoordinatorMessageTypeMemoryRES,
				Payload: exec.GetState().Memory,
			})
			if err != nil {
				log.Println(err)
				continue
			}
		case communication.CoordinatorMessageTypeInputREQ:
			errInput := exec.Input(msg.Payload.(string))
			errInputPayload := ""
			if errInput != nil {
				errInputPayload = errInput.Error()
			}
			err := c.coord.Write(&communication.CoordinatorMessage{
				Type:    communication.CoordinatorMessageTypeInputRES,
				Payload: errInputPayload,
			})
			if err != nil {
				log.Println(err)
				continue
			}
		default:
			log.Println("Unknown message type")
		}
	}
}

func (c *Coordinator) Close() {
	c.conn.Close()
}
