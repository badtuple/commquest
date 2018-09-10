package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

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

func startServer() error {
	log.Println("starting server (TODO)")
	log.Println("server started (TODO)")
	return nil
}

func main() {
	err := startServer()
	if err != nil {
		panic(err)
	}

	log.Println("shutting down")
}
