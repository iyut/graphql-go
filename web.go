package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Settings struct {
	General General  `json:"general"`
	DBInfo  []DBInfo `json:"database"`
}

type General struct {
	PrefixURL string `json:"prefix_url"`
}

type DBInfo struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	DBName   string `json:"dbname"`
}

type PostHandler struct {
	db *sql.DB
}

func (h *PostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var ID int

	rows, err := h.db.Query("Select ID FROM wpa_posts")

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {

		err := rows.Scan(&ID)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprintf(w, "ID : "+strconv.Itoa(ID))
	}
}

func main() {

	settings := openJSONFile()

	prefixURL := settings.General.PrefixURL
	dbInfo := settings.DBInfo[0]

	db, err := sql.Open(dbInfo.Name, dbInfo.Username+":"+dbInfo.Password+"@tcp("+dbInfo.Host+":"+dbInfo.Port+")/"+dbInfo.DBName)

	//db, err := sql.Open("mysql", "root:pass123qwe@tcp(127.0.0.1:3306)/wp_administrator")

	if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("%v", dbInfo)
	}

	defer db.Close()

	r := mux.NewRouter()

	r.PathPrefix(prefixURL + "/").Handler(&PostHandler{db: db})

	http.ListenAndServe(":9990", r)
}

func openJSONFile() Settings {

	jsonFile, err := os.Open("/root/go/bin/settings.json")

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var settings Settings

	json.Unmarshal(byteValue, &settings)

	/*
		if err := json.Unmarshal([]byte(settings), &val); err != nil {
			panic(err)
		}
	*/

	return settings
}
