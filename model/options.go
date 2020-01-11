package model

import "github.com/graph-gophers/graphql-go"

type Options struct {
	OptionID    graphql.ID
	OptionName  string
	OptionValue string
	Autoload    string
}
