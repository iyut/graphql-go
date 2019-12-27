package main

import (
	"context"
	"encoding/json"
	"fmt"

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
		postsByUser(userID: ID!): [Post!]!
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

// Define the greet: String! query:
func (*RootResolver) Greet() string {

	return "Hello World!"
}

func (*RootResolver) GreetPerson(args struct{ Person string }) string {
	return fmt.Sprintf("Hello, %s!", args.Person)
}

type PersonTimeOfDayArgs struct {
	Person    string
	TimeOfDay string
}

var TimesOfDay = map[string]string{
	"MORNING":   "Good morning",
	"AFTERNOON": "Good afternoon",
	"EVENING":   "Good evening",
}

func (*RootResolver) GreetPersonTimeOfDay(ctx context.Context, args PersonTimeOfDayArgs) string {

	timeOfDay, ok := TimesOfDay[args.TimeOfDay]
	if !ok {
		timeOfDay = "Go to bed"
	}

	return fmt.Sprintf("%s %s!", timeOfDay, args.Person)
}

// There are two ways we can define a schema:
//
// - graphql.MustParseSchema(...) *graphql.Schema // Panics on error.
// - graphql.ParseSchema(...) (*graphql.Schema, error)
//
// Define a schema:
var Schema = graphql.MustParseSchema(schemaString, &RootResolver{})

func main() {

	ctx := context.Background()

	///// EXAMPLE 1
	query := `{
		greet	
	}`
	//
	// You can also use these syntax forms if you prefer:
	//
	// descriptiveQuery := `query {
	// 	greet
	// }`
	//
	// moreDescriptiveQuery := `query Greet {
	// 	greet
	// }`
	resp := Schema.Exec(ctx, query, "", nil)
	json0, err := json.MarshalIndent(resp, "", "\t")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(json0))

	///// EXAMPLE 2
	type ClientQuery struct {
		OpName    string                 // Operation name.
		Query     string                 // Query string.
		Variables map[string]interface{} // Query variables (untyped).
	}

	q1 := ClientQuery{
		OpName: "Greet",
		Query: `query Greet {
			greet	
		}`,
		Variables: nil,
	}

	resp1 := Schema.Exec(ctx, q1.Query, q1.OpName, q1.Variables)
	json1, err := json.MarshalIndent(resp1, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(json1))

	///// EXAMPLE 3
	q2 := ClientQuery{
		OpName: "GreetPerson",
		Query: `query GreetPerson($person: String!) {
			greetPerson(person: $person)
		}`,
		Variables: map[string]interface{}{
			"person": "Luthfi",
		},
	}

	resp2 := Schema.Exec(ctx, q2.Query, q2.OpName, q2.Variables)
	json2, err := json.MarshalIndent(resp2, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(json2))

	///// EXAMPLE 4
	q3 := ClientQuery{
		OpName: "GreetPersonTimeOfDay",
		Query: `query GreetPersonTimeOfDay($person: String!, $timeOfDay: TimeOfDay!) {
			greetPersonTimeOfDay( person: $person, timeOfDay: $timeOfDay)
		}`,
		Variables: map[string]interface{}{
			"person":    "Luthfi",
			"timeOfDay": "MORNING",
		},
	}

	resp3 := Schema.Exec(ctx, q3.Query, q3.OpName, q3.Variables)
	json3, err := json.MarshalIndent(resp3, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(json3))
}
