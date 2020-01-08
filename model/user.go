package model

import "github.com/graph-gophers/graphql-go"

type User struct {
	UserID   graphql.ID
	Username string
	Email    string
	Posts    []*Post
}
