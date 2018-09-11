package api

import (
	"database/sql"
	"net/http"

	"../../models"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry

func init() {
	log = logrus.WithField("component", "frontend")
}

type APIFrontend struct{}

func (_ APIFrontend) Name() string {
	return "api"
}

func (_ APIFrontend) Serve() error {
	router := httprouter.New()

	router.POST("/register", wrap(registerPlayerHandler))
	router.GET("/account/:handle", wrap(getPlayerHandler))

	// For preflight options calls
	router.HandleMethodNotAllowed = false

	return http.ListenAndServe(":8080", server{router})
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

func registerPlayerHandler(ctx context) {
	var resp struct {
		Handle string `json:"handle"`
		Name   string `json:"name"`
		Class  string `json:"class"`
	}

	err := ctx.BodyStruct(&resp)
	if err != nil {
		ctx.ServerError(err.Error())
		return
	}

	p, err := models.FindPlayerByHandle(resp.Handle)
	if err == nil {
		ctx.BadRequest("player alreadye exists by that handle")
		return
	}

	if err != sql.ErrNoRows {
		ctx.ServerError(err.Error())
		return
	}

	p, err = models.CreatePlayer(resp.Handle, resp.Name, resp.Class)
	if err != nil {
		ctx.ServerError(err.Error())
		return
	}

	ctx.Ok(p)
}

func getPlayerHandler(ctx context) {
	handle := ctx.Params.ByName("handle")
	p, err := models.FindPlayerByHandle(handle)
	if err != nil {
		if err != sql.ErrNoRows {
			ctx.ServerError(err.Error())
			return
		}

		ctx.BadRequest("player already exists by that handle")
		return
	}

	ctx.Ok(p)
}
