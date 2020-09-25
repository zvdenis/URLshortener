package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"math/big"
	"net/http"
	"strings"
)

type URL struct {
	Address string `json:"address"`
}

//Содержит всю логику для работы с ссылками
type LinkController struct {
	db *sql.DB
}

//Обработчик запросов
type Handler struct {
	linkController *LinkController
}

//Основание
const base int = 62

//Смещение для ID (чтобы ссылки не были слишком короткие)
const bias int64 = 90

func main() {
	db, err := sql.Open("mysql", "root:root@/url_shortener")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var linkController LinkController
	linkController.db = db

	var handler Handler
	handler.linkController = &linkController

	r := mux.NewRouter()

	r.HandleFunc("/", handler.shortenURL).Methods("POST")
	r.HandleFunc("/{link}", handler.redirect).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", r))
}

//Перенаправляет на полную ссылку, используя короткую
func (handler Handler) redirect(w http.ResponseWriter, r *http.Request) {
	//link := mux.Vars(r)["link"]
	fmt.Println(r.URL.Path)
	link := strings.ReplaceAll(r.URL.Path, "/", "")
	fmt.Println(link)
	longLink := handler.linkController.getLongURL(link)
	http.Redirect(w, r, longLink, http.StatusSeeOther)
}

//Сохраняет длинную ссылку и сопоставляет ей короткую
func (handler Handler) shortenURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var longURL URL
	_ = json.NewDecoder(r.Body).Decode(&longURL)
	var shortURL URL
	shortURL.Address = *handler.linkController.genShortURL(&longURL.Address)
	json.NewEncoder(w).Encode(shortURL)
}

//Создает короткую ссылку и сохраняет ее в базе, сопоставляя длинной
func (linkController LinkController) genShortURL(longURL *string) *string {
	shortURL := linkController.getNextShortURL()
	_, err := linkController.db.Exec("insert into url_shortener.links (short_link, long_link) values (?, ?)",
		shortURL, longURL)
	if err != nil {
		panic(err)
	}
	return shortURL
}

//Создает следующую короткую ссылку
func (linkController LinkController) getNextShortURL() *string {
	shortURL := big.NewInt(linkController.getMaxID() + bias + 1).Text(base)
	return &shortURL
}

//Возвращает длинную ссылку по короткой
func (linkController LinkController) getLongURL(shortURL string) string {
	rows, err := linkController.db.Query("SELECT (long_link) FROM url_shortener.links WHERE (short_link = ?);", shortURL)
	if err != nil {
		panic(err)
	}
	var longURL string
	rows.Next()
	rows.Scan(&longURL)
	return longURL
}

//Возвращает максимальный ID в базе
func (linkController LinkController) getMaxID() int64 {
	var id int64 = 0
	rows, err := linkController.db.Query("SELECT MAX(id) FROM url_shortener.links;")
	if err != nil {
		panic(err)
	}
	rows.Next()
	rows.Scan(&id)
	return id
}
