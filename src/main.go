package main

import (
	"time"

	"github.com/sirupsen/logrus"

	"./config"
	"./frontends"
	"./game"
	"./util"
)

var log *logrus.Entry = util.LoggerFor("main")

func init() {
	log.Println("initializing Commquest")
	log.Println("loading config")
	err := config.Load()
	if err != nil {
		panic(err)
	}
	log.Printf("loaded [%s] config", config.Get().Name)
}

func startFrontendAPI() error {
	log.Println("starting frontend server")
	fe := frontends.New(config.Get().Frontend)

	err := fe.Serve()
	if err != nil {
		return err
	}

	return nil
}

func startGameLoop() {
	for range time.Tick(10 * time.Second) {
		game.PlayTurn()
	}
}

func main() {
	// TODO: Graceful shutdown

	go startGameLoop()

	err := startFrontendAPI()
	if err != nil {
		log.Fatal("frontend server couldn't be started: %v", err.Error())
	}

	log.Println("shutting down")
}
