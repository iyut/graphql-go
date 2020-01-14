package model

import (
	"time"

	"github.com/graph-gophers/graphql-go"
)

type User struct {
	UserID            graphql.ID
	Username          string
	Email             string
	UserLogin         string
	UserPassword      string
	UserNicename      string
	UserEmail         string
	UserURL           string
	UserRegistered    time.Time
	UserActivationKey string
	UserStatus        int32
	DisplayName       string
	Posts             []*Post
	UserMeta          []*UserMeta
}

type UserMeta struct {
	UMetaID   graphql.ID
	UserID    graphql.ID
	MetaKey   string
	MetaValue string
}
