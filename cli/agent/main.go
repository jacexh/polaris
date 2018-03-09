package main

import "github.com/jacexh/polaris/agent"

func main() {
	task, err := agent.NewSniffTask("10.0.1.36", 8000, 0, agent.ConsolePrinter{}.Handle)
	if err != nil {
		panic(err)
	}
	task.Sniff()

}
