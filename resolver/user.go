package resolver

import (
	"database/sql"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/iyut/graphql-go/model"
)

/*
 * UserResolver
 *
 * type User {
 * 	userID: ID!
 * 	username: String!
 * 	emoji: String!
 * 	notes: [Note!]!
 * }
 */

type UserResolver struct {
	U  *model.User
	DB *sql.DB
}

func (r *UserResolver) UserID() graphql.ID {
	return r.U.UserID
}

func (r *UserResolver) Username() string {
	return r.U.Username
}

func (r *UserResolver) Email() string {
	return r.U.Email
}

func (r *UserResolver) Posts() ([]*PostResolver, error) {
	rootRxs := &RootResolver{DB: r.DB}

	return rootRxs.Posts(struct{ UserID graphql.ID }{UserID: r.U.UserID})
}
