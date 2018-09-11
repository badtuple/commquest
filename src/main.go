package main

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/sirupsen/logrus"

	"./frontends"
	"./game"
	"./util"
)

var log *logrus.Entry = util.LoggerFor("main")

func init() {
	log.Println("initializing Commquest")
	err := loadConfig()
	if err != nil {
		panic(err)
	}
}

var config struct {
	// Name of the specific configuration.
	// Useful if running different quests for different channels/services
	Name string `json:"name"`

	// The frontend service used for the game.
	// Right now "api" is the only option, but will eventually have things
	// like "slack", "irc", "discord", "twitch", etc.
	Frontend string `json:"frontend"`
}

func loadConfig() error {
	log.Println("loading config")
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		return err
	}

	log.Printf("loaded [%s] config", config.Name)
	return nil
}

func startFrontendAPI() error {
	log.Println("starting frontend server")

	fe := frontends.New(config.Frontend)
	err := fe.Serve()
	if err != nil {
		return err
	}

	log.Println("server started")
	return nil
}

func startGameLoop() {
	for range time.Tick(1 * time.Minute) {
		game.PlayTurn()
	}
}

func main() {
	// TODO: Graceful shutdown

	go startGameLoop()

	err := startFrontendAPI()
	if err != nil {
		panic(err)
	}

	log.Println("shutting down")
}
