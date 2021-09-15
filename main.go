package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path"

	gocache "github.com/mrod502/go-cache"
	"github.com/mrod502/logger"
)

func main() {
	var b []byte
	var err error

	cfgPath := flag.String("config-path", "serverConfig.json", "the path to the config file")
	home, _ := os.UserHomeDir()
	flag.Parse()

	b, err = ioutil.ReadFile(path.Join(home, *cfgPath))
	if err != nil {
		pwd, _ := os.Getwd()
		b, err = ioutil.ReadFile(path.Join(pwd, *cfgPath))
		if err != nil {
			panic(err)
		}
	}

	var config logger.ServerConfig
	config.KeySignatures = gocache.NewBoolCache()
	err = json.Unmarshal(b, &config)

	if err != nil {
		panic(err)
	}
	var server logger.LogServer

	server, err = logger.NewLogServer(config)
	if err != nil {
		panic(err)
	}
	defer server.Quit()
	go func() {
		err := server.Serve()
		if err != nil {
			panic(err)
		}
	}()

	logger.Info("LOGGER", "listening on port", fmt.Sprintf("%d", config.Port))
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
	fmt.Println("bye")

}
