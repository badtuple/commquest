package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type APIFrontend struct{}

func (_ APIFrontend) Name() string {
	return "api"
}

func (_ APIFrontend) Serve() error {
	router := httprouter.New()

	router.POST("/register", wrap(registerPlayerHandler))
	router.GET("/account/:id", wrap(getPlayerHandler))

	// For preflight options calls
	router.HandleMethodNotAllowed = false

	return http.ListenAndServe(":8080", server{router})
}

type context struct {
	Response http.ResponseWriter
	Request  *http.Request
	Params   *httprouter.Params
	//Player *models.Player
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

func wrap(h func(context)) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := context{w, r, &p}
		ctx.Response.Header().Set("Content-Type", "application/json")

		// TODO: auth

		h(ctx)
	}
}

func registerPlayerHandler(ctx context) {}
func getPlayerHandler(ctx context)      {}
