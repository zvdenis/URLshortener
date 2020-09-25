package main

import (
	"URLshortener/Links"
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)

type URL struct {
	Address string `json:"address"`
}

//Обработчик запросов
type Handler struct {
	linkController *Links.LinkController
}

func main() {
	db, err := sql.Open("mysql", "root:root@/url_shortener")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var linkController Links.LinkController
	linkController.Db = db

	var handler Handler
	handler.linkController = &linkController

	r := mux.NewRouter()

	r.HandleFunc("/", handler.shortenURL).Methods("POST")
	r.HandleFunc("/{link}", handler.redirect).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", r))
}

//Перенаправляет на полную ссылку, используя короткую
func (handler Handler) redirect(w http.ResponseWriter, r *http.Request) {
	link := strings.ReplaceAll(r.URL.Path, "/", "")
	longLink := handler.linkController.GetLongURL(link)
	http.Redirect(w, r, longLink, http.StatusSeeOther)
}

//Сохраняет длинную ссылку и сопоставляет ей короткую
func (handler Handler) shortenURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var longURL URL
	_ = json.NewDecoder(r.Body).Decode(&longURL)
	var shortURL URL
	shortURL.Address = *handler.linkController.GenShortURL(&longURL.Address)
	json.NewEncoder(w).Encode(shortURL)
}
