package main

import "github.com/jacexh/polaris/agent"

func main() {
	task, err := agent.NewSniffTask("127.0.0.1", 8000, 0, agent.ConsolePrinter{}.Handle)
	if err != nil {
		panic(err)
	}
	task.Sniff()
}
