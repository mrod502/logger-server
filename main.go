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

const revision = "0.0.1"

func version() string {
	return fmt.Sprintf("logger-server version %s\nlogger version %s", revision, logger.Version())
}

func main() {
	var b []byte
	var err error

	cfgPath := flag.String("config-path", "serverConfig.json", "the path to the config file")
	versionFlag := flag.Bool("version", false, "print current version")
	generateApiKey := flag.Bool("gen-api-key", false, "generate a new api key")
	generateKeyPair := flag.Bool("gen-key-pair", false, "generate a new pub/private key pair")
	home, _ := os.UserHomeDir()
	flag.Parse()

	if *versionFlag {
		fmt.Println(version())
		return
	}

	if *generateApiKey {
		key, sig := logger.GenerateApiKey()
		fmt.Println("api key:", key)
		fmt.Println("sha256sum:", sig)
		return
	}
	if *generateKeyPair {
		key, err := logger.GenerateKey(4096)
		if err != nil {
			panic(err)
		}
		priv, pub, err := logger.MarshalKeyPair(key)

		if err != nil {
			panic(err)
		}
		fmt.Println(string(priv))
		fmt.Println(string(pub))
		return
	}

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
