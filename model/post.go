package model

import (
	"time"

	"github.com/graph-gophers/graphql-go"
)

type Post struct {
	PostID              graphql.ID
	PostAuthor          graphql.ID
	PostDate            time.Time
	PostDateGMT         time.Time
	PostContent         string
	PostTitle           string
	PostExcerpt         string
	PostStatus          string
	CommentStatus       string
	PingStatus          string
	PostPassword        string
	PostName            string
	ToPing              string
	Pinged              string
	PostModified        time.Time
	PostModifiedGMT     time.Time
	PostContentFiltered string
	PostParent          graphql.ID
	GUID                string
	MenuOrder           int32
	PostType            string
	PostMimeType        string
	CommentCount        int64
	PostMeta            []*PostMeta
	Author              *User
	TermRelationships   []*TermRelationships
}

type PostMeta struct {
	MetaID    graphql.ID
	PostID    graphql.ID
	MetaKey   string
	MetaValue string
}

type PostInput struct {
	Title string
}
