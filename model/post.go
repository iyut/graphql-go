package model

import "github.com/graph-gophers/graphql-go"

type Post struct {
	PostID graphql.ID
	Title  string
}

type PostInput struct {
	Title string
}
