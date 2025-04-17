package router

import (
	"database/sql"
	"net/http"
)

func WebRouter(db *sql.DB) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
}
