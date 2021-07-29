package endpoint

import (
	"errors"
	"log"
	"net"
	"steel-lang/semantics"
	"steel-simulator-common/communication"
	"steel-simulator-common/config"
)

// AgentEndpoint wraps a communication endpoint to add agent functionality
type AgentEndpoint struct {
	end *communication.Endpoint
}

// New creates a new endpoint, connected to the coordinator
func New() (*AgentEndpoint, error) {
	// I resolve the address for the coordinator...
	tcpAddr, err := net.ResolveTCPAddr("tcp", "steel-coordinator:5001")
	if err != nil {
		return nil, err
	}
	// ... I connect to it...
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}
	// ... and I return the endpoint
	return &AgentEndpoint{
		end: communication.New(conn),
	}, nil
}

// SendInit sends the initialization message to the coordinator
func (a *AgentEndpoint) SendInit(name string) error {
	// I write the INIT message, sending the agent name...
	err := a.end.Write(&communication.EndpointMessage{
		Type:    communication.EndpointMessageTypeINIT,
		Payload: name,
	})
	if err != nil {
		return err
	}
	// ... and I read the ACK
	msg, err := a.end.Read()
	if err != nil {
		return err
	}
	if msg.Type != communication.EndpointMessageTypeACK {
		return errors.New("unexpected response to init")
	}
	return nil
}

// HandleMessages listens for messages and responds to them
func (a *AgentEndpoint) HandleMessages(exec *semantics.MuSteelExecuter, agent config.Agent) {
	for {
		// I read a message...
		msg, err := a.end.Read()
		if err != nil {
			log.Println(err)
			break
		}
		// ... and I check its type
		switch msg.Type {
		// If it is a memory request...
		case communication.EndpointMessageTypeMemoryREQ:
			// ... I respond with the memory state
			err := a.end.Write(&communication.EndpointMessage{
				Type:    communication.EndpointMessageTypeMemoryRES,
				Payload: exec.GetState().Memory,
			})
			if err != nil {
				log.Println(err)
				continue
			}
		// If it is a command to input...
		case communication.EndpointMessageTypeInputREQ:
			// ... I execute it...
			errInput := exec.Input(msg.Payload.(string))
			errInputPayload := ""
			if errInput != nil {
				errInputPayload = errInput.Error()
			}
			// ... and I respond with the eventual error
			err := a.end.Write(&communication.EndpointMessage{
				Type:    communication.EndpointMessageTypeInputRES,
				Payload: errInputPayload,
			})
			if err != nil {
				log.Println(err)
				continue
			}
		// If it is a configuration request...
		case communication.EndpointMessageTypeConfigREQ:
			// ... I respond with the initialization configuration
			err := a.end.Write(&communication.EndpointMessage{
				Type:    communication.EndpointMessageTypeConfigRES,
				Payload: agent,
			})
			if err != nil {
				log.Println(err)
				continue
			}
		// Otherwise I cannot do anything
		default:
			log.Println("Unknown message type")
		}
	}
}

// Close closes the endpoint connection
func (e *AgentEndpoint) Close() {
	e.end.Close()
}
