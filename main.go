package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	graphql "github.com/graph-gophers/graphql-go"
)

/****
*********************
GET THE SETTINGS INFO
*********************
****/
type Settings struct {
	General General  `json:"general"`
	DBInfo  []DBInfo `json:"database"`
}

type General struct {
	PrefixURL     string `json:"prefix_url"`
	GraphqlSchema string `json:"graphql_schema"`
}

type DBInfo struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	DBName   string `json:"dbname"`
}

/****
*********************
SET THE TYPE FOR GRAPHQL
*********************
****/
type User struct {
	UserID   graphql.ID
	Username string
	Email    string
	Posts    []*Post
}

type Post struct {
	PostID graphql.ID
	Title  string
}

type PostInput struct {
	Title string
}

// Define mock data:
var users = []*User{
	{
		UserID:   graphql.ID("u-001"),
		Username: "nyxerys",
		Email:    "nyxerys@nyxerys.com",
		Posts: []*Post{
			{PostID: "n-001", Title: "Olá Mundo!"},
			{PostID: "n-002", Title: "Olá novamente, mundo!"},
			{PostID: "n-003", Title: "Olá, escuridão!"},
		},
	}, {
		UserID:   graphql.ID("u-002"),
		Username: "rdnkta",
		Email:    "rdnkta@rdnkta.com",
		Posts: []*Post{
			{PostID: "n-004", Title: "Привіт Світ!"},
			{PostID: "n-005", Title: "Привіт ще раз, світ!"},
			{PostID: "n-006", Title: "Привіт, темрява!"},
		},
	}, {
		UserID:   graphql.ID("u-003"),
		Username: "username_ZAYDEK",
		Email:    "username_ZAYDEK@zaydek.com",
		Posts: []*Post{
			{PostID: "n-007", Title: "Hello, world!"},
			{PostID: "n-008", Title: "Hello again, world!"},
			{PostID: "n-009", Title: "Hello, darkness!"},
		},
	},
}

/****
*********************
DEFINE THE RESOLVER FOR GRAPHQL
*********************
****/

/*
 * RootResolver
 *
 * type User {
 * 	userID: ID!
 * 	username: String!
 * 	emoji: String!
 * 	notes: [Note!]!
 * }
 */
type RootResolver struct {
	db *sql.DB
}

func (r *RootResolver) Users() ([]*UserResolver, error) {

	var userRxs []*UserResolver

	rows, err := r.db.Query(`
		SELECT
			ID,
			user_login,
			user_email
		FROM
			wpa_users
	`)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {

		user := &User{}
		err := rows.Scan(&user.UserID, &user.Username, &user.Email)

		if err != nil {
			return nil, err
		}

		userRxs = append(userRxs, &UserResolver{user})
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return userRxs, nil
}

func (r *RootResolver) User(args struct{ UserID graphql.ID }) (*UserResolver, error) {

	var useridInt int64
	user := &User{}
	err := r.db.QueryRow(`
		SELECT
			ID,
			user_login,
			user_email
		FROM
			wpa_users
		WHERE
			ID = ?
	`, args.UserID).Scan(&useridInt, &user.Username, &user.Email)

	useridStr := strconv.FormatInt(useridInt, 10)
	user.UserID = graphql.ID(useridStr)
	if err != nil {
		return nil, err
	}

	return &UserResolver{user}, nil
}

func (r *RootResolver) Posts(args struct{ UserID graphql.ID }) ([]*PostResolver, error) {

	var postIDInt int64
	var postRxs []*PostResolver
	rows, err := r.db.Query(`
		SELECT
			ID,
			post_title
		FROM
			wpa_posts
		WHERE
			post_author = ?
	`, args.UserID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		post := &Post{}
		err := rows.Scan(&postIDInt, &post.Title)

		if err != nil {
			return nil, err
		}

		postIDStr := strconv.FormatInt(postIDInt, 10)
		post.PostID = graphql.ID(postIDStr)

		postRxs = append(postRxs, &PostResolver{post})
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return postRxs, nil
}

func (r *RootResolver) Post(args struct{ PostID graphql.ID }) (*PostResolver, error) {

	var postIDInt int64
	post := &Post{}
	err := r.db.QueryRow(`
		SELECT
			ID,
			post_title
		FROM
			wpa_posts
		WHERE
			ID = ?
	`, args.PostID).Scan(&postIDInt, &post.Title)

	if err != nil {
		return nil, err
	}

	postIDStr := strconv.FormatInt(postIDInt, 10)
	post.PostID = graphql.ID(postIDStr)

	return &PostResolver{post}, nil

}

type CreatePostArgs struct {
	UserID graphql.ID
	Post   PostInput
}

func (r *RootResolver) CreatePost(args CreatePostArgs) (*PostResolver, error) {

	var post *Post

	for _, user := range users {
		// Create a note with a note ID of n-010:
		post = &Post{PostID: "n-010", Title: args.Post.Title}
		user.Posts = append(user.Posts, post) // Push note.
	}
	// Return note:
	return &PostResolver{post}, nil

}

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
	u *User
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
	rootRxs := &RootResolver{}

	return rootRxs.Posts(struct{ UserID graphql.ID }{UserID: r.u.UserID})
}

/*
 * PostResolver
 *
 * type Post {
 * 	PostID: ID!
 * 	title: String!
 * }
 */

type PostResolver struct {
	p *Post
}

func (r *PostResolver) PostID() graphql.ID {
	return r.p.PostID
}

func (r *PostResolver) Title() string {
	return r.p.Title
}

func main() {

	ctx := context.Background()

	settings := openJSONFile()

	//prefixURL := settings.General.PrefixURL
	dbInfo := settings.DBInfo[0]

	db, err := sql.Open(dbInfo.Name, dbInfo.Username+":"+dbInfo.Password+"@tcp("+dbInfo.Host+":"+dbInfo.Port+")/"+dbInfo.DBName)
	//db, err := sql.Open("mysql", "root:pass123qwe@tcp(127.0.0.1:3306)/wp_administrator")

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	bstr, err := ioutil.ReadFile(settings.General.GraphqlSchema)
	if err != nil {
		panic(err)
	}

	schemaString := string(bstr)

	schema, err := graphql.ParseSchema(schemaString, &RootResolver{db: db})
	if err != nil {
		panic(err)
	}

	type JSON = map[string]interface{}

	type ClientQuery struct {
		OpName    string
		Query     string
		Variables JSON
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

	resp1 := schema.Exec(ctx, q1.Query, q1.OpName, q1.Variables)
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
			"userID": "1",
		},
	}

	resp2 := schema.Exec(ctx, q2.Query, q2.OpName, q2.Variables)
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
			"userID": "1",
		},
	}

	resp3 := schema.Exec(ctx, q3.Query, q3.OpName, q3.Variables)
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
			"postID": "1",
		},
	}

	resp4 := schema.Exec(ctx, q4.Query, q4.OpName, q4.Variables)
	json4, err := json.MarshalIndent(resp4, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(json4))

	q5 := ClientQuery{
		OpName: "CreatePost",
		Query: `mutation CreatePost($userID: ID!, $post: PostInput!){
				createPost(userID: $userID, post: $post){
					postID,
					title
				}
			}`,
		Variables: JSON{
			"userID": "u-0003",
			"post": JSON{
				"title": "We create a post!",
			},
		},
	}

	resp5 := schema.Exec(ctx, q5.Query, q5.OpName, q5.Variables)
	json5, err := json.MarshalIndent(resp5, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(json5))

	q6 := ClientQuery{
		OpName: "Users",
		Query: `query Users{
				users{
					userID
					username
					email
					posts {
						postID
						title
					}
				}
			}`,
		Variables: nil,
	}

	resp6 := schema.Exec(ctx, q6.Query, q6.OpName, q6.Variables)
	json6, err := json.MarshalIndent(resp6, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(json6))
}

func openJSONFile() Settings {

	jsonFile, err := os.Open("/root/go/bin/settings.json")

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var settings Settings

	json.Unmarshal(byteValue, &settings)

	/*
		if err := json.Unmarshal([]byte(settings), &val); err != nil {
			panic(err)
		}
	*/

	return settings
}
