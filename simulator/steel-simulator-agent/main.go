package main

import (
	"os"
	"steel-lang/communication"
	"steel-lang/semantics"
	"steel-simulator-agent/endpoint"
	"steel-simulator-agent/memory"
	"steel-simulator-common/config"
	"time"

	"log"
)

func main() {
	// I check if a config is present on the Args...
	if len(os.Args) < 2 {
		log.Fatalln("Config not found, exiting")
	}
	// ... and I deserialize it to get its fields
	configStr := os.Args[1]
	agent := config.Agent{}
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
	// ... and I create the executer
	log.Println("Creating executer")
	exec, err := semantics.NewMuSteelExecuter(mem, agent.Rules, communication.MakeMemberlistAgent(mem.ResourceNames(), 5000, agent.Endpoints))
	if err != nil {
		log.Fatal(err)
	}
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
	go end.HandleMessages(exec, agent)
	// Finally, I start the executer loop
	log.Println("Starting main loop")
	for {
		// I execute a command...
		exec.Exec()
		// ... and I sleep for a while
		time.Sleep(agent.Tick)
	}
}
