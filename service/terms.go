package service

import (
	"database/sql"

	"github.com/graph-gophers/graphql-go"
	"github.com/iyut/graphql-go/helper"
	"github.com/iyut/graphql-go/model"
)

func NewTermsService(db *sql.DB, prefix string) *Terms {

	return &Terms{db: db, prefix: prefix}
}

type Terms struct {
	db     *sql.DB
	prefix string
}

type ArgsTerms struct {
	TermID   int64
	Taxonomy string
	Slug     string
	ParentID int64
}

func (t *Terms) GetTerms(args ArgsTerms) ([]*model.TermTaxonomy, error) {

	var terms []*model.TermTaxonomy

	var termTaxID int64
	var termID int64
	var termGroup int64
	var termParent int64

	var queryMap []interface{}

	query := `
	SELECT
		tt.term_taxonomy_id,
		tt.term_id,
		tt.taxonomy,
		tt.description,
		tt.parent,
		tt.count,
		t.name,
		t.slug,
		t.term_group,
	FROM
	` + t.prefix + "term_taxonomy tt, " + t.prefix + "terms t" + `
	WHERE
		tt.term_id = t.term_id
	`

	if args.TermID > 0 {
		query = query + " AND tt.term_id = ? "
		queryMap = append(queryMap, args.TermID)
	}

	if len(args.Taxonomy) > 0 {
		query = query + " AND tt.taxonomy = ? "
		queryMap = append(queryMap, args.Taxonomy)
	}

	if len(args.Slug) > 0 {
		query = query + " AND t.slug = ? "
		queryMap = append(queryMap, args.Slug)
	}

	if args.ParentID > 0 {
		query = query + " AND tt.parent = ? "
		queryMap = append(queryMap, args.ParentID)
	}

	query = query + ";"

	rows, err := t.db.Query(query, queryMap...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		termTax := &model.TermTaxonomy{}
		term := &model.Terms{}

		err := rows.Scan(&termTaxID, &termID, &termTax.Taxonomy, &termTax.Description, &termParent, &termTax.Count, &term.Name, &term.Slug, &termGroup)

		termTax.TermTaxonomyID = helper.IntToGraphqlID(termTaxID)
		termTax.TermID = helper.IntToGraphqlID(termID)
		termTax.Parent = helper.IntToGraphqlID(termParent)

		term.TermID = termTax.TermID
		term.TermGroup = helper.IntToGraphqlID(termGroup)

		if err != nil {
			return nil, err
		}

		termsMeta, err := t.GetTermsMeta(termTax.TermID)

		if err == nil {
			term.TermsMeta = termsMeta
		}

		termTax.Terms = term

		terms = append(terms, termTax)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return terms, err
}

func (t *Terms) GetTermsMeta(termID graphql.ID) ([]*model.TermsMeta, error) {

	var termMetas []*model.TermsMeta
	var tMetaIDInt int64
	var termIDInt int64

	rows, err := t.db.Query(`
		SELECT
			meta_id,
			term_id,
			meta_key,
			meta_value
		FROM
	`+t.prefix+"term_meta"+`
		WHERE
			term_id = ?
	`, termID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		meta := &model.TermsMeta{}
		err := rows.Scan(&tMetaIDInt, &termIDInt, &meta.MetaKey, &meta.MetaValue)

		if err != nil {
			return nil, err
		}

		meta.MetaID = helper.IntToGraphqlID(tMetaIDInt)
		meta.TermID = helper.IntToGraphqlID(termIDInt)

		termMetas = append(termMetas, meta)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return termMetas, nil
}

func (t *Terms) TaxonomyExist(taxonomy string) bool {

	var count int

	row := t.db.QueryRow(`
		SELECT
			count( taxonomy )
		FROM
	`+t.prefix+"term_taxonomy"+`
		WHERE
			taxonomy = ?
	`, taxonomy)

	err := row.Scan(&count)
	if err != nil {
		return false
	}

	if count > 0 {
		return true
	}

	return false
}
