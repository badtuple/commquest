package slack_frontend

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

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
	var ev slackEvent
	err := json.NewDecoder(r.Body).Decode(&ev)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch ev.Event.Type {
	case "member_joined_channel":
		log.Printf("member_joined_channel: %+v", ev)
	case "member_left_channel":
		log.Printf("member_left_channel: %+v", ev)
	default:
		log.Printf("unexpected event type: %+v", ev)
	}
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
