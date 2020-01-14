package service

import "database/sql"

func NewPostService(db *sql.DB, prefix string) *Post {

	return &Post{db: db, prefix: prefix}
}

type Post struct {
	db     *sql.DB
	prefix string
}
