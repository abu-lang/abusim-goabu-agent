package main

import (
	"os"
	"steel-lang/communication"
	"steel-lang/datastructure"
	"steel-lang/semantics"
	"steel-simulator-agent/coordinator"
	"steel-simulator-agent/memory"
	"steel-simulator-config/config"
	"sync"
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
	coord, err := coordinator.New()
	if err != nil {
		log.Fatalln(err)
	}
	defer coord.Close()
	err = coord.SendInit(configStr)
	if err != nil {
		log.Fatalln(err)
	}
	mem, err := memory.New(agent.MemoryController, agent.Memory)
	if err != nil {
		log.Fatalln(err)
	}
	for {
		// log.Println(agent)
		log.Println(mem.GetResources())
		time.Sleep(5 * time.Second)
	}
}

func Example() {
	memory := datastructure.MakeResources()
	memory.Bool["button"] = false
	memory.Bool["involved"] = false
	r1 := "rule r1 on button; default involved = false; for all ext.involved == false do involved = true;"
	rules := []string{r1}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		e1, err := semantics.NewMuSteelExecuter(memory, rules, communication.MakeMemberlistAgent(memory.ResourceNames(), 10001, nil))
		if err != nil {
			log.Fatal(err)
		}
		time.AfterFunc(2500*time.Millisecond, func() {
			e1.Input("button = true;")
		})
		for {
			e1.Exec()
			mem1 := e1.GetState().Memory
			if !mem1.Bool["involved"] {
				log.Println("1: involved should be true")
			} else {
				break
			}
			time.Sleep(1 * time.Second)
		}
		log.Println("OK1")
	}()

	go func() {
		defer wg.Done()
		e2, err := semantics.NewMuSteelExecuter(memory, rules, communication.MakeMemberlistAgent(memory.ResourceNames(), 10002, []string{"localhost:10001"}))
		if err != nil {
			log.Fatal(err)
		}
		for {
			e2.Exec()
			mem1 := e2.GetState().Memory
			if !mem1.Bool["involved"] {
				log.Println("2: involved should be true")
			} else {
				break
			}
			time.Sleep(1 * time.Second)
		}
		log.Println("OK2")
	}()

	wg.Wait()
}
