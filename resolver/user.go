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

	if len(r.U.Username) > 0 {
		return r.U.Username
	} else {
		return r.U.UserLogin
	}

}

func (r *UserResolver) Email() string {

	if len(r.U.Email) > 0 {
		return r.U.Email
	} else {
		return r.U.UserEmail
	}

}

func (r *UserResolver) Nicename() string {
	return r.U.UserNicename
}

func (r *UserResolver) Status() int32 {
	return r.U.UserStatus
}

func (r *UserResolver) Posts() ([]*PostResolver, error) {
	rootRxs := &RootResolver{DB: r.DB}

	return rootRxs.Posts(struct{ UserID graphql.ID }{UserID: r.U.UserID})
}
