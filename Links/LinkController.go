package Links

import (
	"database/sql"
	"math/big"
)

//Основание
const base int = 62

//Смещение для ID (чтобы ссылки не были слишком короткие)
const bias int64 = 90

//Содержит всю логику для работы с ссылками
type LinkController struct {
	Db *sql.DB
}

//Создает короткую ссылку и сохраняет ее в базе, сопоставляя длинной
func (linkController LinkController) GenShortURL(longURL *string) *string {
	shortURL := linkController.getNextShortURL()
	_, err := linkController.Db.Exec("insert into url_shortener.links (short_link, long_link) values (?, ?)",
		shortURL, longURL)
	if err != nil {
		panic(err)
	}
	return shortURL
}

//Создает следующую короткую ссылку
func (linkController LinkController) getNextShortURL() *string {
	shortURL := big.NewInt(linkController.GetMaxID() + bias + 1).Text(base)
	return &shortURL
}

//Возвращает длинную ссылку по короткой
func (linkController LinkController) GetLongURL(shortURL string) string {
	rows, err := linkController.Db.Query("SELECT (long_link) FROM url_shortener.links WHERE (short_link = ?);", shortURL)
	if err != nil {
		panic(err)
	}
	var longURL string
	rows.Next()
	rows.Scan(&longURL)
	return longURL
}

//Возвращает максимальный ID в базе
func (linkController LinkController) GetMaxID() int64 {
	var id int64 = 0
	rows, err := linkController.Db.Query("SELECT MAX(id) FROM url_shortener.links;")
	if err != nil {
		panic(err)
	}
	rows.Next()
	rows.Scan(&id)
	return id
}
