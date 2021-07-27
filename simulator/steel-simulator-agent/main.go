package main

import (
	"os"
	"steel-lang/communication"
	"steel-lang/semantics"
	"steel-simulator-agent/coordinator"
	"steel-simulator-agent/memory"
	"steel-simulator-common/config"
	"time"

	"log"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Config not found, exiting")
	}
	configStr := os.Args[1]
	agent := config.Agent{}
	err := agent.Deserialize(configStr)
	if err != nil {
		log.Fatalf("Bad config deserialization: %v", err)
	}
	log.Println("Creating memory")
	mem, err := memory.New(agent.MemoryController, agent.Memory)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Creating executer")
	exec, err := semantics.NewMuSteelExecuter(mem, agent.Rules, communication.MakeMemberlistAgent(mem.ResourceNames(), 5000, agent.Endpoints))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connecting to coordinator")
	coord, err := coordinator.New()
	if err != nil {
		log.Fatalln(err)
	}
	defer coord.Close()
	err = coord.SendSelfName(agent.Name)
	if err != nil {
		log.Fatalln(err)
	}
	go coord.HandleMessages(exec)
	log.Println("Starting main loop")
	for {
		exec.Exec()
		time.Sleep(agent.Tick)
	}
}
