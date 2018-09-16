package frontend

import (
	"errors"

	"../util"
	//"./api"
	slack "./slack"
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry = util.LoggerFor("frnt")

// Global state is bad, but I'm basically treating this
// package as more of a module. Eventually this will likely
// need to be a slice of frontends if we are going to
// support running the same game on multiple at once
var fe Frontend

// All frontends must implement this interface.
type Frontend interface {
	// Name of the Frontend (api, irc, twitch, etc...)
	Name() string

	Serve() error

	PushMessage(string) error
}

// expected to block
func Serve(name string) error {
	log.Printf("starting [%v] frontend", name)

	switch name {
	case "api":
	// TODO: We broke the API fe while moving the slack
	// api to RTM.  Will need to hook it back up

	//fe = Frontend(api.APIFrontend{})
	case "slack":
		fe = Frontend(slack.SlackFrontend{})
		fe.Serve()
	default:
		log.Fatalf("unimplemented frontend %s", name)
	}
	return nil
}

// Push a message to the running frontend
func PushMessage(msg string) error {
	if fe == nil {
		return errors.New("frontend not initialized")
	}

	return fe.PushMessage(msg)

}
