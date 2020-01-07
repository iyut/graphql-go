package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	graphql "github.com/graph-gophers/graphql-go"
)

/****
*********************
RESPONSE STATUS URL
*********************
****/
const (
	StatusCodeOK              = 200
	StatusCodeBadRequest      = 400
	StatusCodeUnauthorized    = 401
	StatusCodeRequestFailed   = 402
	StatusCodeNotFound        = 404
	StatusCodeConflict        = 409
	StatusCodeTooManyRequests = 429
	StatusCodeServerError     = 500
)

var Statuses = map[int]string{
	StatusCodeOK:              "OK",
	StatusCodeBadRequest:      "Bad Request",
	StatusCodeUnauthorized:    "Unauthorized",
	StatusCodeRequestFailed:   "Request Failed",
	StatusCodeNotFound:        "Not Found",
	StatusCodeConflict:        "Conflict",
	StatusCodeTooManyRequests: "Too Many Requests",
	StatusCodeServerError:     "Server Error",
}

var (
	RespondOK              = NewResponder(StatusCodeOK)
	RespondBadRequest      = NewResponder(StatusCodeBadRequest)
	RespondUnauthorized    = NewResponder(StatusCodeUnauthorized)
	RespondRequestFailed   = NewResponder(StatusCodeRequestFailed)
	RespondNotFound        = NewResponder(StatusCodeNotFound)
	RespondConflict        = NewResponder(StatusCodeConflict)
	RespondTooManyRequests = NewResponder(StatusCodeTooManyRequests)
	RespondServerError     = NewResponder(StatusCodeServerError)
)

func NewResponder(statusCode int) func(http.ResponseWriter) {
	respond := func(w http.ResponseWriter) {
		if statusCode >= 200 && statusCode <= 299 {
			w.WriteHeader(statusCode)
			return
		}
		status := Statuses[statusCode]
		http.Error(w, status, statusCode)
	}
	return respond
}

type GraphqlHandler struct {
	Schema *graphql.Schema
}

func (h *GraphqlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		RespondNotFound(w)
		return
	}

	if r.Body == nil {
		RespondServerError(w)
		log.Printf("Request Body: %s", "No query data")
		return
	}

	type JSON = map[string]interface{}

	type ClientQuery struct {
		OperationName string `json:"operationName"`
		Query         string `json:"query"`
		Variables     JSON   `json:"variables"`
	}

	var rBody ClientQuery

	err := json.NewDecoder(r.Body).Decode(&rBody)

	if err != nil {
		RespondServerError(w)
		log.Printf("Error parsing JSON request body")
		return
	}

	q1 := ClientQuery{
		OperationName: rBody.OperationName,
		Query:         rBody.Query,
		Variables:     rBody.Variables,
	}

	resp1 := h.Schema.Exec(context.Background(), q1.Query, q1.OperationName, q1.Variables)
	if len(resp1.Errors) > 0 {
		RespondServerError(w)
		log.Printf("Schema.Exec: %+v", resp1.Errors)
		return
	}

	json1, err := json.MarshalIndent(resp1, "", "\t")
	if err != nil {
		RespondServerError(w)
		log.Printf("json.MarshalIndent: %s", err)
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	fmt.Fprintf(w, string(json1))
}
