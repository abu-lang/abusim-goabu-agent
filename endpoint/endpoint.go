package endpoint

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/abu-lang/abusim-core/schema"

	"github.com/abu-lang/goabu"
	steelconfig "github.com/abu-lang/goabu/config"
)

// nameToLogLevel converts from a log level name to the corresponding level
var nameToLogLevel = map[string]int{
	"Fatal":   steelconfig.LogFatal,
	"Error":   steelconfig.LogError,
	"Warning": steelconfig.LogWarning,
	"Info":    steelconfig.LogInfo,
	"Debug":   steelconfig.LogDebug,
}

// logLevelToName converts from a log level to the corresponding name
var logLevelToName = map[int]string{
	steelconfig.LogFatal:   "Fatal",
	steelconfig.LogError:   "Error",
	steelconfig.LogWarning: "Warning",
	steelconfig.LogInfo:    "Info",
	steelconfig.LogDebug:   "Debug",
}

// AgentEndpoint wraps a schema endpoint to add agent functionality
type AgentEndpoint struct {
	end *schema.Endpoint
}

// New creates a new endpoint, connected to the coordinator
func New() (*AgentEndpoint, error) {
	// I resolve the address for the coordinator...
	tcpAddr, err := net.ResolveTCPAddr("tcp", "abusim-coordinator:5001")
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
		end: schema.New(conn),
	}, nil
}

// SendInit sends the initialization message to the coordinator
func (a *AgentEndpoint) SendInit(name string) error {
	// I write the INIT message, sending the agent name...
	payload := schema.EndpointMessagePayloadINIT{
		Name: name,
	}
	err := a.end.Write(&schema.EndpointMessage{
		Type:    schema.EndpointMessageTypeINIT,
		Payload: &payload,
	})
	if err != nil {
		return err
	}
	// ... and I read the ACK
	msg, err := a.end.Read()
	if err != nil {
		return err
	}
	if msg.Type != schema.EndpointMessageTypeACK {
		return errors.New("unexpected response to init")
	}
	return nil
}

// HandleMessages listens for messages and responds to them
func (a *AgentEndpoint) HandleMessages(exec *goabu.Executer, agent schema.AgentConfiguration, paused *bool) {
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
		case schema.EndpointMessageTypeMemoryREQ:
			// ... I get the state...
			state := exec.TakeState()
			// ... I get the memory from the state...
			memory := schema.MemoryResources{}
			memory.Bool = state.Memory.Bool
			memory.Integer = state.Memory.Integer
			memory.Float = state.Memory.Float
			memory.Text = state.Memory.Text
			memory.Time = state.Memory.Time
			// ... I get a string representation of the pool...
			pool := [][]schema.PoolElem{}
			for _, ruleActions := range state.Pool {
				poolActions := []schema.PoolElem{}
				for _, action := range ruleActions {
					poolActions = append(poolActions, schema.PoolElem{
						Resource: action.Resource,
						Value:    fmt.Sprintf("%v", action.Value),
					})
				}
				pool = append(pool, poolActions)
			}
			// ... and I respond with the state
			payload := schema.EndpointMessagePayloadMemoryRES{
				Memory: memory,
				Pool:   pool,
			}
			err := a.end.Write(&schema.EndpointMessage{
				Type:    schema.EndpointMessageTypeMemoryRES,
				Payload: &payload,
			})
			if err != nil {
				log.Println(err)
				continue
			}
		// If it is a command to input...
		case schema.EndpointMessageTypeInputREQ:
			// ... I execute it...
			errInput := exec.Input(msg.Payload.(*schema.EndpointMessagePayloadInputREQ).Input)
			errInputPayload := ""
			if errInput != nil {
				errInputPayload = errInput.Error()
			}
			// ... and I respond with the eventual error
			payload := schema.EndpointMessagePayloadInputRES{
				Error: errInputPayload,
			}
			err := a.end.Write(&schema.EndpointMessage{
				Type:    schema.EndpointMessageTypeInputRES,
				Payload: &payload,
			})
			if err != nil {
				log.Println(err)
				continue
			}
		// If it is a debug status request...
		case schema.EndpointMessageTypeDebugREQ:
			// ... I respond with the debug status
			payload := schema.EndpointMessagePayloadDebugRES{
				Paused:    *paused,
				Verbosity: logLevelToName[exec.LogLevel()],
			}
			err := a.end.Write(&schema.EndpointMessage{
				Type:    schema.EndpointMessageTypeDebugRES,
				Payload: &payload,
			})
			if err != nil {
				log.Println(err)
				continue
			}
		// If it is a debug status change...
		case schema.EndpointMessageTypeDebugChangeREQ:
			// ... I execute it...
			newStatus := msg.Payload.(*schema.EndpointMessagePayloadDebugChangeREQ)
			*paused = newStatus.Paused
			exec.SetLogLevel(nameToLogLevel[newStatus.Verbosity])
			// ... and I respond
			err := a.end.Write(&schema.EndpointMessage{
				Type:    schema.EndpointMessageTypeDebugChangeRES,
				Payload: nil,
			})
			if err != nil {
				log.Println(err)
				continue
			}
		// If it is a debug step request...
		case schema.EndpointMessageTypeDebugStepREQ:
			// ... I step the executer...
			exec.Exec()
			// ... and I respond
			err := a.end.Write(&schema.EndpointMessage{
				Type:    schema.EndpointMessageTypeDebugStepRES,
				Payload: nil,
			})
			if err != nil {
				log.Println(err)
				continue
			}
		// If it is a configuration request...
		case schema.EndpointMessageTypeConfigREQ:
			// ... I respond with the initialization configuration
			payload := schema.EndpointMessagePayloadConfigRES{
				Agent: agent,
			}
			err := a.end.Write(&schema.EndpointMessage{
				Type:    schema.EndpointMessageTypeConfigRES,
				Payload: &payload,
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
