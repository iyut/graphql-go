package main

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	graphql "github.com/graph-gophers/graphql-go"
)

// Define a schema string:
const schemaString = `
	# Define what the schema is capable of:
	schema {
		query: Query
	}

	type User {
		userID: ID!
		username: String!
		email: String!
		posts: [Post!]!
	}

	type Post{
		postID: ID!
		title: String!
	}

	# Define what the queries are capable of:
	type Query {
		# List Users:
		users: [User!]!
		# Get User:
		user(userID: ID!): User!
		#List Post per User:
		posts(userID: ID!): [Post!]!
		#Get Post:
		post(postID: ID!): Post!
	}
`

type User struct {
	UserID   graphql.ID
	Username string
	Email    string
	Posts    []Post
}

type Post struct {
	PostID graphql.ID
	Title  string
}

// Define mock data:
var users = []User{
	{
		UserID:   graphql.ID("u-001"),
		Username: "nyxerys",
		Email:    "nyxerys@nyxerys.com",
		Posts: []Post{
			{PostID: "n-001", Title: "Olá Mundo!"},
			{PostID: "n-002", Title: "Olá novamente, mundo!"},
			{PostID: "n-003", Title: "Olá, escuridão!"},
		},
	}, {
		UserID:   graphql.ID("u-002"),
		Username: "rdnkta",
		Email:    "rdnkta@rdnkta.com",
		Posts: []Post{
			{PostID: "n-004", Title: "Привіт Світ!"},
			{PostID: "n-005", Title: "Привіт ще раз, світ!"},
			{PostID: "n-006", Title: "Привіт, темрява!"},
		},
	}, {
		UserID:   graphql.ID("u-003"),
		Username: "username_ZAYDEK",
		Email:    "username_ZAYDEK@zaydek.com",
		Posts: []Post{
			{PostID: "n-007", Title: "Hello, world!"},
			{PostID: "n-008", Title: "Hello again, world!"},
			{PostID: "n-009", Title: "Hello, darkness!"},
		},
	},
}

// Define a root resolver to hook queries onto:
type RootResolver struct{}

func (r *RootResolver) Users() ([]User, error) {

	return users, nil
}

func (r *RootResolver) User(args struct{ UserID graphql.ID }) (User, error) {

	for _, user := range users {
		if args.UserID == user.UserID {
			return user, nil
		}
	}

	return User{}, nil
}

func (r *RootResolver) Posts(args struct{ UserID graphql.ID }) ([]Post, error) {

	user, err := r.User(args)
	if reflect.ValueOf(user).IsZero() || err != nil {

		return nil, err
	}

	return user.Posts, nil
}

func (r *RootResolver) Post(args struct{ PostID graphql.ID }) (Post, error) {

	for _, user := range users {
		for _, post := range user.Posts {

			if args.PostID == post.PostID {
				return post, nil
			}
		}
	}

	return Post{}, nil
}

// There are two ways we can define a schema:
//
// - graphql.MustParseSchema(...) *graphql.Schema // Panics on error.
// - graphql.ParseSchema(...) (*graphql.Schema, error)
//
// Define a schema:
var opts = []graphql.SchemaOpt{graphql.UseFieldResolvers()}
var Schema = graphql.MustParseSchema(schemaString, &RootResolver{}, opts...)

func main() {

	ctx := context.Background()

	type ClientQuery struct {
		OpName    string
		Query     string
		Variables map[string]interface{}
	}

	q1 := ClientQuery{
		OpName: "Users",
		Query: `query Users{
			users{
				userID
				username
				email
			}
		}`,
		Variables: nil,
	}

	resp1 := Schema.Exec(ctx, q1.Query, q1.OpName, q1.Variables)
	json1, err := json.MarshalIndent(resp1, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(json1))

	q2 := ClientQuery{
		OpName: "User",
		Query: `query User($userID: ID!){
			user(userID: $userID){
				userID
				username
				email
			}
		}`,
		Variables: map[string]interface{}{
			"userID": "u-001",
		},
	}

	resp2 := Schema.Exec(ctx, q2.Query, q2.OpName, q2.Variables)
	json2, err := json.MarshalIndent(resp2, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(json2))

	q3 := ClientQuery{
		OpName: "Posts",
		Query: `query Posts($userID: ID!){
			posts(userID: $userID){
				postID
				title
			}
		}`,
		Variables: map[string]interface{}{
			"userID": "u-002",
		},
	}

	resp3 := Schema.Exec(ctx, q3.Query, q3.OpName, q3.Variables)
	json3, err := json.MarshalIndent(resp3, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(json3))

	q4 := ClientQuery{
		OpName: "Post",
		Query: `query Post($postID: ID!){
			post(postID: $postID){
				postID
				title
			}
		}`,
		Variables: map[string]interface{}{
			"postID": "n-007",
		},
	}

	resp4 := Schema.Exec(ctx, q4.Query, q4.OpName, q4.Variables)
	json4, err := json.MarshalIndent(resp4, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(json4))
}
