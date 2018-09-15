package frontend

import (
	"errors"
	"net/http"

	"../util"
	"./api"
	slack "./slack"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry = util.LoggerFor("frnt")

// Global state is bad.  But I'm basically treating this
// package as more of a module.  Eventually this wil likely
// need to be a slice of frontends if we are going to
// support running the same game on multiple at once
var fe Frontend

// All frontends must implement this interface.
type Frontend interface {
	// Name of the Frontend (api, irc, twitch, etc...)
	Name() string

	// Get router to serve
	Router() *httprouter.Router
	// Get Port to serve on
	Port() string

	PushMessage(string) error
}

func Serve(name string) error {
	log.Printf("starting [%v] frontend", name)

	switch name {
	case "api":
		fe = Frontend(api.APIFrontend{})
	case "slack":
		fe = Frontend(slack.SlackFrontend{})
	default:
		log.Fatalf("unimplemented frontend %s", name)
	}

	router := fe.Router()
	return http.ListenAndServe(fe.Port(), server{router})
}

type server struct {
	r *httprouter.Router
}

func (s server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, HEAD, OPTIONS, PUT, PATCH, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	s.r.ServeHTTP(w, r)
}

// Push a message to the running frontend
func PushMessage(msg string) error {
	if fe == nil {
		return errors.New("frontend not initialized")
	}

	return fe.PushMessage(msg)

}
