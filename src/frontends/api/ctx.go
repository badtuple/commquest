package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type context struct {
	Response http.ResponseWriter
	Request  *http.Request
	Params   *httprouter.Params
}

func (ctx *context) Ok(data interface{}) {
	resp, err := json.Marshal(data)
	if err != nil {
		log.Println(err.Error())
		http.Error(ctx.Response, "500 error marshalling data", http.StatusInternalServerError)
		return
	}
	log.Println(string(resp))
	fmt.Fprint(ctx.Response, bytes.NewBuffer(resp))
}

type ErrorResponse struct {
	Errors []string `json:"errors"`
}

func (e ErrorResponse) toJSON() string {
	resp, _ := json.Marshal(e)
	return string(resp)
}

func newErrorResponse(e string) string {
	resp, _ := json.Marshal(ErrorResponse{[]string{e}})
	return string(resp)
}

func (ctx *context) ServerError(e string) {
	http.Error(ctx.Response, newErrorResponse(e), 500)
}

func (ctx *context) BadRequest(e string) {
	http.Error(ctx.Response, newErrorResponse(e), 400)
}

func (ctx *context) Unauthorized(e string) {
	http.Error(ctx.Response, newErrorResponse(e), 401)
}

func (ctx *context) NotFound() {
	http.Error(ctx.Response, newErrorResponse("not found"), 404)
}

func (ctx *context) Body() (b map[string]interface{}, err error) {
	return b, json.NewDecoder(ctx.Request.Body).Decode(&b)
}

func (ctx *context) BodyStruct(target interface{}) error {
	return json.NewDecoder(ctx.Request.Body).Decode(target)
}
