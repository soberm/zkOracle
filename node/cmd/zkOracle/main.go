package main

import (
	"flag"
	"github.com/spf13/viper"
	"net"
	"node/pkg/zkOracle"
	"os"
	"os/signal"
)

func main() {

	configFile := flag.String("c", "./configs/config.example.json", "filename of the config file")
	flag.Parse()

	viper.SetConfigFile(*configFile)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	var config zkOracle.Config
	err := viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}

	node, err := zkOracle.NewNode(&config)
	if err != nil {
		panic(err)
	}

	listener, err := net.Listen("tcp", config.BindAddress)
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
