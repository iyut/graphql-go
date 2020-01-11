package model

import (
	"time"

	"github.com/graph-gophers/graphql-go"
)

type Links struct {
	LinkID          graphql.ID
	LinkURL         string
	LinkName        string
	LinkImage       string
	LinkTarget      string
	LinkDescription string
	LinkVisible     string
	LinkOwner       int64
	LinkRating      int32
	LinkUpdate      time.Time
	LinkRel         string
	LinkNotes       string
	LinkRSS         string
}
