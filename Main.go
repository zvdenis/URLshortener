package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"math/big"
	"net/http"
)

type URL struct {
	Address string `json:"address"`
}

var database *sql.DB

const base int = 62

func main() {
	db, err := sql.Open("mysql", "root:root@/url_shortener")
	database = db

	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/", shortenURL).Methods("POST")
	r.HandleFunc("/{link}", redirect).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", r))
}

func redirect(w http.ResponseWriter, r *http.Request) {
	link := mux.Vars(r)["link"]
	longLink := getLongURL(link)
	http.Redirect(w, r, longLink, 301)
}

func shortenURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var longURL URL
	_ = json.NewDecoder(r.Body).Decode(&longURL)
	shortURL := getNextShortURL()

	_, err := database.Exec("insert into url_shortener.links (short_link, long_link) values (?, ?)",
		shortURL.Address, longURL.Address)
	if err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(shortURL)
}

func getNextShortURL() URL {
	var shortURL URL
	shortURL.Address = big.NewInt(getMaxID() + 1).Text(base)
	return shortURL
}

func getLongURL(shortURL string) string {
	rows, err := database.Query("SELECT (long_link) FROM url_shortener.links WHERE (short_link = ?);", shortURL)
	if err != nil {
		panic(err)
	}
	var longURL string
	rows.Next()
	rows.Scan(&longURL)
	return longURL
}

func getMaxID() int64 {
	var id int64 = 0
	rows, err := database.Query("SELECT MAX(id) FROM url_shortener.links;")
	if err != nil {
		panic(err)
	}
	rows.Next()
	rows.Scan(&id)
	return id
}
