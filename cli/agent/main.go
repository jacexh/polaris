package main

import "github.com/jacexh/polaris/agent"

func main() {
	client, err := agent.NewWSClient("ws://127.0.0.1:16666/ws/register")
	if err != nil {
		panic(err)
	}
	if err = client.Register(); err != nil {
		panic(err)
	}
	task, err := agent.NewSniffTask("127.0.0.1", 8000, 0, agent.ConsolePrinter{}.Handle)
	if err != nil {
		panic(err)
	}
	task.Sniff()
}
