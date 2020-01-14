package service

import (
	"database/sql"
	"errors"
	"time"

	"github.com/graph-gophers/graphql-go"
	"github.com/iyut/graphql-go/helper"
	"github.com/iyut/graphql-go/model"
)

func NewUserService(db *sql.DB, prefix string) *User {

	return &User{db: db, prefix: prefix}
}

type User struct {
	db     *sql.DB
	prefix string
}

type ArgsUser struct {
	UserID   int64
	Email    string
	Slug     string
	Username string
}

func (u *User) GetUsers(args ArgsUser) ([]*model.User, error) {

	var users []*model.User

	var useridInt int64
	var userRegisteredString string

	var queryMap []interface{}

	query := `
		SELECT
			ID,
			user_login,
			user_pass,
			user_nicename,
			user_email,
			user_url,
			user_registered,
			user_activation_key,
			user_status,
			display_name
		FROM	
	` + u.prefix + "users" + `
		WHERE
			1 = 1
	`

	if args.UserID > 0 {
		query = query + " AND ID = ? "
		queryMap = append(queryMap, args.UserID)
	}

	if len(args.Email) > 0 {
		query = query + " AND user_email = ? "
		queryMap = append(queryMap, args.Email)
	}

	if len(args.Slug) > 0 {
		query = query + " AND user_nicename = ? "
		queryMap = append(queryMap, args.Slug)
	}

	if len(args.Username) > 0 {
		query = query + " AND user_login = ? "
		queryMap = append(queryMap, args.Username)
	}

	query = query + ";"

	rows, err := u.db.Query(query, queryMap...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		user := &model.User{}

		err := rows.Scan(&useridInt, &user.UserLogin, &user.UserPassword, &user.UserNicename, &user.UserEmail, &user.UserURL, &userRegisteredString, &user.UserActivationKey, &user.UserStatus, &user.DisplayName)

		if err != nil {
			return nil, err
		}

		user.UserID = helper.IntToGraphqlID(useridInt)
		user.UserRegistered, err = time.Parse("2006-01-02 15:04:05", userRegisteredString)

		if err != nil {
			return nil, err
		}

		userMeta, err := u.GetMeta(user.UserID)

		if err == nil {
			user.UserMeta = userMeta
		}

		users = append(users, user)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return users, err
}

func (u *User) FindBy(field string, value string) (*model.User, error) {

	acceptedFields := [4]string{"id", "email", "slug", "username"}
	usedFields := [4]string{"ID", "user_email", "user_nicename", "user_login"}
	var usedField string

	if helper.ItemExists(acceptedFields, field) == false {
		return nil, errors.New("field is not accepted")
	}

	for i := 0; i < len(acceptedFields); i++ {
		if acceptedFields[i] == field {
			usedField = usedFields[i]
		}
	}

	var useridInt int64
	user := &model.User{}

	err := u.db.QueryRow(`
		SELECT
			ID,
			user_login,
			user_pass,
			user_nicename,
			user_email,
			user_url,
			user_registered,
			user_activation_key,
			user_status,
			display_name
		
		FROM	
	`+u.prefix+"users"+`
		WHERE
			1 = 1
			AND `+usedField+` = ?
	`, value).Scan(
		&useridInt,
		&user.UserLogin,
		&user.UserPassword,
		&user.UserNicename,
		&user.UserEmail,
		&user.UserURL,
		&user.UserRegistered,
		&user.UserActivationKey,
		&user.UserStatus,
		&user.DisplayName)

	user.UserID = helper.IntToGraphqlID(useridInt)

	userMeta, err := u.GetMeta(user.UserID)
	if err == nil {
		user.UserMeta = userMeta
	}

	return user, nil
}

func (u *User) GetMeta(userID graphql.ID) ([]*model.UserMeta, error) {

	var userMetas []*model.UserMeta
	var uMetaIDInt int64
	var userIDInt int64

	rows, err := u.db.Query(`
		SELECT
			umeta_id,
			user_id,
			meta_key,
			meta_value
		FROM
	`+u.prefix+"user_meta"+`
		WHERE
			user_id = ?
	`, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		userMeta := &model.UserMeta{}
		err := rows.Scan(&uMetaIDInt, &userIDInt, &userMeta.MetaKey, &userMeta.MetaValue)

		if err != nil {
			return nil, err
		}

		userMeta.UMetaID = helper.IntToGraphqlID(uMetaIDInt)
		userMeta.UserID = helper.IntToGraphqlID(userIDInt)

		userMetas = append(userMetas, userMeta)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return userMetas, nil
}

func (u *User) GetPosts(userID graphql.ID) ([]*model.Post, error) {

	var posts []*model.Post

	return posts, nil
}
