package slack_frontend

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"../../config"
	"../../models"
	"../../util"

	"github.com/julienschmidt/httprouter"
	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry = util.LoggerFor("frnt")

var api = slack.New(config.Get().Credentials.SlackAPIKey)

type SlackFrontend struct{}

func (_ SlackFrontend) Name() string {
	return "slack"
}

func (_ SlackFrontend) Router() *httprouter.Router {
	router := httprouter.New()
	router.POST("/events", eventHandler)

	// For preflight options calls
	router.HandleMethodNotAllowed = false

	return router
}

func (_ SlackFrontend) Port() string {
	return ":8081"
}

func eventHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var ev slackEvent
	err := json.NewDecoder(r.Body).Decode(&ev)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch ev.Event.Type {
	case "member_joined_channel":
		maybeCreateUser(ev)
	case "member_left_channel":
		// Right now do nothing. When a User being
		// active is impemented then we'd set them
		// to inactive.
		log.Printf("member_left_channel: %+v", ev)
	default:
		log.Printf("unexpected event type: %+v", ev)
	}
}
func maybeCreateUser(ev slackEvent) {
	uid := ev.Event.User
	log.Printf("member_joined_channel: %+v", uid)
	user, err := api.GetUserInfo(uid)
	if err != nil {
		log.Printf("could not get user %v: %v", uid, err.Error())
		return
	}

	p, err := models.FindPlayerByHandle(user.Name)
	if err == nil {
		log.Printf("user by handle %v already exists", user.Name)
		return
	}

	if err != sql.ErrNoRows {
		log.Println("%+v", err.Error())
		return
	}

	title := user.Profile.Title
	p, err = models.CreatePlayer(user.Name, user.Name, title)
	if err != nil {
		log.Println("could not create player: %v", err.Error())
		return
	}

	log.Printf("created player %+v", p)
}

type slackEvent struct {
	Type    string `json:"type"`
	EventID string `json:"event_id"`

	Event struct {
		// "member_joined_channel" or "member_left_channel"
		Type    string  `json:"type"`
		EventTS float64 `json:"event_ts"`
		User    string  `json:"user"`
		Channel string  `json:"channel"`
	} `json:"event"`
}
