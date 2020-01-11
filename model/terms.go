package model

import "github.com/graph-gophers/graphql-go"

type Terms struct {
	TermID    graphql.ID
	Name      string
	Slug      string
	TermGroup graphql.ID
	TermsMeta []*TermsMeta
}

type TermsMeta struct {
	MetaID    graphql.ID
	TermID    graphql.ID
	MetaKey   string
	MetaValue string
}

type TermTaxonomy struct {
	TermTaxonomyID graphql.ID
	TermID         graphql.ID
	Taxonomy       string
	Description    string
	Parent         graphql.ID
	Count          uint64
	Terms          *Terms
}

type TermRelationships struct {
	ObjectID       graphql.ID
	TermTaxonomyID graphql.ID
	TermOrder      int64
	TermTaxonomy   *TermTaxonomy
	Object         *Post
}
