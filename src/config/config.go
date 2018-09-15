package config

import (
	"encoding/json"
	"io/ioutil"
)

var config Config

// The main configuration struct read in fron config.json
// All configuration should happen in here and nowhere else.
type Config struct {
	// Name of the specific configuration.
	// Useful if running different quests for different channels/services
	Name string `json:"name"`

	// The frontend service used for the game.
	// Right now "api" is the only option, but will eventually have things
	// like "slack", "irc", "discord", "twitch", etc.
	Frontend string `json:"frontend"`

	// Credentials for connecting to frontends
	Credentials struct {
		SlackAPIKey string `json:"slack_api_key"`
	} `json:"credentials"`
}

func Load() error {
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		return err
	}
	return nil
}

func Get() Config {
	return config
}
