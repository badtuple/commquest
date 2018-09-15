package api

import (
	"database/sql"
	"net/http"

	"../../models"
	"../../util"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry = util.LoggerFor("frnt")

type APIFrontend struct{}

func (_ APIFrontend) Name() string {
	return "api"
}

func (_ APIFrontend) Router() *httprouter.Router {
	router := httprouter.New()

	router.POST("/register", wrap(registerPlayerHandler))
	router.GET("/account/:handle", wrap(getPlayerHandler))

	// For preflight options calls
	router.HandleMethodNotAllowed = false
	return router
}

func (_ APIFrontend) Port() string {
	return ":8080"
}

// The API frontend doesn't have the concept of pushing a
// message.  Evetually we may want to implement a webhook
// style system where the Config can contain a URL we'll
// send to.
func (_ APIFrontend) PushMessage(msg string) error {
	return nil
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
