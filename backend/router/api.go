package router

import (
	"database/sql"
	"net/http"

	"real_time_forum/backend/auth"
)

func ApiRouter(db *sql.DB) {
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		auth.Login(w, r, db)
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		auth.Register(w, r, db)
	})
}
