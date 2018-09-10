package frontends

import (
	"log"

	"./api"
)

// All frontends must implement this interface.
type Frontend interface {
	// Name of the Frontend (api, irc, twitch, etc...)
	Name() string
	// Start the server listening for updates from chat client.
	Serve() error
}

func New(name string) Frontend {
	switch name {
	case "api":
		return Frontend(api.APIFrontend{})
	default:
		log.Fatalf("unimplemented frontend %s", name)
	}

	return nil
}