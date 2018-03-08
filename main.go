package main

import "github.com/jacexh/polaris/agent"

func main() {
	sniffer, err := agent.NewSniffer("10.0.1.36", 8000)
	if err != nil {
		panic(err)
	}
	sniffer.Run()
}
