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
	var resp map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	raw, ok := resp["type"]
	if !ok {
		log.Printf("event has no type: %+v", resp)
		return
	}

	typ, ok := raw.(string)
	if !ok {
		log.Printf("event type is not a string: %+v", resp)
		return
	}

	switch typ {
	case "member_joined_channel":
		log.Printf("member_joined_channel: %+v", resp)
	case "member_left_channel":
		log.Printf("member_left_channel: %+v", resp)
	}
}
