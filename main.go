package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/iyut/graphql-go/handler"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
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
	GraphqlURL    string `json:"graphql_url"`
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

		userRxs = append(userRxs, &UserResolver{u: user, db: r.db})
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

	return &UserResolver{u: user, db: r.db}, nil
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

		postRxs = append(postRxs, &PostResolver{p: post, db: r.db})
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

	return &PostResolver{p: post, db: r.db}, nil

}

type CreatePostArgs struct {
	UserID graphql.ID
	Post   PostInput
}

func (r *RootResolver) CreatePost(args CreatePostArgs) (*PostResolver, error) {

	res, err := r.db.Exec(`
		INSERT INTO wpa_posts (
			post_author,
			post_title )
		VALUES (?, ?);
	`, args.UserID, args.Post.Title)

	if err != nil {
		return nil, err
	}

	lastid, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	postIDStr := strconv.FormatInt(lastid, 10)

	return r.Post(struct{ PostID graphql.ID }{PostID: graphql.ID(postIDStr)})

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

/*
 * PostResolver
 *
 * type Post {
 * 	PostID: ID!
 * 	title: String!
 * }
 */

type PostResolver struct {
	p  *Post
	db *sql.DB
}

func (r *PostResolver) PostID() graphql.ID {
	return r.p.PostID
}

func (r *PostResolver) Title() string {
	return r.p.Title
}

func main() {

	settings := openJSONFile()

	graphqlURL := settings.General.GraphqlURL
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

	//params := r.URL.Query()
	schema, err := graphql.ParseSchema(schemaString, &RootResolver{db: db})
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.PathPrefix(graphqlURL).Handler(&handler.GraphqlHandler{Schema: schema})
	r.PathPrefix(graphqlURL + "/").Handler(&handler.GraphqlHandler{Schema: schema})

	http.ListenAndServe(":9990", r)
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
