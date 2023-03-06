package main

import (
	"flag"
	"github.com/spf13/viper"
	"node/pkg/zkOracle"
	"os"
	"os/signal"
)

func main() {

	configFile := flag.String("c", "./configs/config.example.json", "filename of the config file")
	flag.Parse()

	v := viper.GetViper()
	err := zkOracle.LoadConfig(v, *configFile)
	if err != nil {
		panic(err)
	}

	node, err := zkOracle.NewNode(v)
	if err != nil {
		panic(err)
	}

	go func() {
		if err := node.Run(); err != nil {
			panic(err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	node.Stop()
}
