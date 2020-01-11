package service

import (
	"database/sql"
	"errors"

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

func (u *User) FindBy(field string, value string) (*model.User, error) {

	acceptedField := [4]string{"id", "email", "slug", "login"}
	if helper.ItemExists(acceptedField, field) == false {
		return nil, errors.New("field is not accepted")
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
	`+u.prefix+"_users"+`
		WHERE
			ID = ?
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
	`+u.prefix+"_user_meta"+`
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
