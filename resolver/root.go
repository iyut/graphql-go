package resolver

import (
	"database/sql"
	"strconv"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/iyut/graphql-go/model"
	"github.com/iyut/graphql-go/service"
)

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
	DB *sql.DB
}

func (r *RootResolver) Users() ([]*UserResolver, error) {

	var userRxs []*UserResolver

	userService := service.NewUserService(r.DB, "wpa_")

	argsUser := service.ArgsUser{}
	users, err := userService.GetUsers(argsUser)

	if err != nil {
		return nil, err
	}

	for _, user := range users {
		userRxs = append(userRxs, &UserResolver{U: user, DB: r.DB})
	}
	/*
		rows, err := r.DB.Query(`
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

			user := &model.User{}
			err := rows.Scan(&user.UserID, &user.Username, &user.Email)

			if err != nil {
				return nil, err
			}

			userRxs = append(userRxs, &UserResolver{U: user, DB: r.DB})
		}

		err = rows.Err()
		if err != nil {
			return nil, err
		}
	*/

	return userRxs, nil
}

func (r *RootResolver) User(args struct{ UserID graphql.ID }) (*UserResolver, error) {

	var useridInt int64
	user := &model.User{}
	err := r.DB.QueryRow(`
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

	return &UserResolver{U: user, DB: r.DB}, nil
}

func (r *RootResolver) Posts(args struct{ UserID graphql.ID }) ([]*PostResolver, error) {

	var postIDInt int64
	var postRxs []*PostResolver
	rows, err := r.DB.Query(`
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

		post := &model.Post{}
		err := rows.Scan(&postIDInt, &post.PostTitle)

		if err != nil {
			return nil, err
		}

		postIDStr := strconv.FormatInt(postIDInt, 10)
		post.PostID = graphql.ID(postIDStr)

		postRxs = append(postRxs, &PostResolver{P: post, DB: r.DB})
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return postRxs, nil
}

func (r *RootResolver) Post(args struct{ PostID graphql.ID }) (*PostResolver, error) {

	var postIDInt int64
	post := &model.Post{}
	err := r.DB.QueryRow(`
		SELECT
			ID,
			post_title
		FROM
			wpa_posts
		WHERE
			ID = ?
	`, args.PostID).Scan(&postIDInt, &post.PostTitle)

	if err != nil {
		return nil, err
	}

	postIDStr := strconv.FormatInt(postIDInt, 10)
	post.PostID = graphql.ID(postIDStr)

	return &PostResolver{P: post, DB: r.DB}, nil

}

type CreatePostArgs struct {
	UserID graphql.ID
	Post   model.PostInput
}

func (r *RootResolver) CreatePost(args CreatePostArgs) (*PostResolver, error) {

	res, err := r.DB.Exec(`
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
