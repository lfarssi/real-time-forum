package router

import (
	"database/sql"
	"net/http"

	"real_time_forum/backend/auth"
	"real_time_forum/backend/controllers"
)

func ApiRouter(db *sql.DB) {
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		auth.Login(w, r, db)
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		auth.Register(w, r, db)
	})

	http.HandleFunc("/addPost", func(w http.ResponseWriter, r *http.Request) {
		controllers.AddPost(w, r, db)
	})
}
