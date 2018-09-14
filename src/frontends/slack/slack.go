package slack_frontend

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
)

// TODO: Replace with config option. Don't push to git
var TOKEN = "TOKEN"
var api = slack.New(TOKEN)

type SlackFrontend struct{}

func (_ SlackFrontend) Name() string {
	return "slack"
}

func (_ SlackFrontend) Serve() error {
	router := httprouter.New()
	router.POST("/events", eventHandler)
	router.HandleMethodNotAllowed = false // For preflight options calls

	log.Println("started [slack] frontend")
	return http.ListenAndServe(":8081", server{router})
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

func eventHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	body := buf.String()
	eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{TOKEN}))
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	log.Printf("%+v", body)

	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "text")
		w.Write([]byte(r.Challenge))
	}

	//if eventsAPIEvent.Type == slackevents.CallbackEvent {
	//postParams := slack.PostMessageParameters{}
	//innerEvent := eventsAPIEvent.InnerEvent

	//switch ev := innerEvent.Data.(type) {
	//case *slackevents.AppMentionEvent:
	//api.PostMessage(ev.Channel, "Yes, hello.", postParams)
	//}
	//}

	switch eventsAPIEvent.Type {
	case "member_joined_channel":
		log.Printf("member_joined_channel: %+v", eventsAPIEvent)
	case "member_left_channel":
		log.Printf("member_left_channel: %+v", eventsAPIEvent)
	}
}

//member_joined_channel
//A user joined a public or private channel
//channels:read

//member_left_channel
//A user left a public or private channel
//channels:read
