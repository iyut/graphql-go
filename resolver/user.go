package resolver

import (
	"database/sql"

	"github.com/graph-gophers/graphql-go"
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
	u  *User
	db *sql.DB
}

func (r *UserResolver) UserID() graphql.ID {
	return r.u.UserID
}

func (r *UserResolver) Username() string {
	return r.u.Username
}

func (r *UserResolver) Email() string {
	return r.u.Email
}

func (r *UserResolver) Posts() ([]*PostResolver, error) {
	rootRxs := &RootResolver{db: r.db}

	return rootRxs.Posts(struct{ UserID graphql.ID }{UserID: r.u.UserID})
}
