package main

import (
	"os"
	"time"

	"github.com/abu-lang/abusim-core/abusim-goabu-agent/endpoint"
	"github.com/abu-lang/abusim-core/abusim-goabu-agent/memory"

	"github.com/abu-lang/abusim-core/schema"
	"github.com/abu-lang/goabu"
	"github.com/abu-lang/goabu/communication"
	goabuconfig "github.com/abu-lang/goabu/config"

	"log"
)

func main() {
	// I check if a config is present on the Args...
	if len(os.Args) < 2 {
		log.Fatalln("Config not found, exiting")
	}
	// ... and I deserialize it to get its fields
	configStr := os.Args[1]
	agent := schema.AgentConfiguration{}
	err := agent.Deserialize(configStr)
	if err != nil {
		log.Fatalf("Bad config deserialization: %v", err)
	}
	// I create the memory for the agent...
	log.Println("Creating memory")
	mem, err := memory.New(agent.MemoryController, agent.Memory)
	if err != nil {
		log.Fatalln(err)
	}
	// ... I create the executer...
	log.Println("Creating executer")
	logConfig := goabuconfig.LogConfig{
		Encoding: "console",
		Level:    goabuconfig.LogError,
	}
	exec, err := goabu.NewExecuter(mem, agent.Rules, communication.NewMemberlistAgent(agent.Name, 5000, logConfig, agent.Endpoints...), logConfig)
	if err != nil {
		log.Fatal(err)
	}
	// ... and I create the paused variable
	paused := false
	// I connect to the coordinator...
	log.Println("Connecting to coordinator")
	end, err := endpoint.New()
	if err != nil {
		log.Fatalln(err)
	}
	defer end.Close()
	// ... I send to it the initialization message...
	err = end.SendInit(agent.Name)
	if err != nil {
		log.Fatalln(err)
	}
	// ... and I start the main message loop
	go end.HandleMessages(exec, agent, &paused)
	// Finally, I start the executer loop
	log.Println("Starting main loop")
	for {
		// I execute a command if not paused...
		if !paused {
			exec.Exec()
		}
		// ... and I sleep for a while
		time.Sleep(agent.Tick)
	}
}
