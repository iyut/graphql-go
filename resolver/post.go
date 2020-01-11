package resolver

import (
	"database/sql"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/iyut/graphql-go/model"
)

/*
 * PostResolver
 *
 * type Post {
 * 	PostID: ID!
 * 	title: String!
 * }
 */

type PostResolver struct {
	P  *model.Post
	DB *sql.DB
}

func (r *PostResolver) PostID() graphql.ID {
	return r.P.PostID
}

func (r *PostResolver) Title() string {
	return r.P.PostTitle
}
