package model

import "github.com/graph-gophers/graphql-go"

import "time"

type Comments struct {
	CommentID          graphql.ID
	CommentPostID      graphql.ID
	CommentAuthor      bool
	CommentAuthorEmail string
	CommentAuthorURL   string
	CommentAuthorIP    string
	CommentDate        time.Time
	CommentDateGMT     time.Time
	CommentContent     string
	CommentKarma       int64
	CommentApproved    string
	CommentAgent       string
	CommentType        string
	CommentParent      graphql.ID
	UserID             graphql.ID
}

type CommentsMeta struct {
	MetaID    graphql.ID
	CommentID graphql.ID
	MetaKey   string
	MetaValue string
}
