package main

import (
	"net"
	"node/pkg/zkOracle"
	"os"
	"os/signal"
)

func main() {

	node, err := zkOracle.NewNode()
	if err != nil {
		panic(err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:25565")
	if err != nil {
		panic(err)
	}

	go func() {
		if err := node.Run(listener); err != nil {
			panic(err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	node.Stop()
}
