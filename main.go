package main

import (
	"database/sql"
	"io/ioutil"
	"net/http"

	"github.com/iyut/graphql-go/handler"
	"github.com/iyut/graphql-go/resolver"

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

	bstr, err := ioutil.ReadFile(getMainPath() + "main-schema.graphql")
	if err != nil {
		panic(err)
	}

	schemaString := string(bstr)

	//params := r.URL.Query()
	schema, err := graphql.ParseSchema(schemaString, &resolver.RootResolver{DB: db})
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.PathPrefix(graphqlURL).Handler(&handler.GraphqlHandler{Schema: schema})
	r.PathPrefix(graphqlURL + "/").Handler(&handler.GraphqlHandler{Schema: schema})

	http.ListenAndServe(":9990", r)
}
