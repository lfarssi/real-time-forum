package router

import (
	"html/template"
	"net/http"

	"real_time_forum/backend/controllers"
)

func WebRouter() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/" {
			http.Error(w, "Page not found", http.StatusNotFound)
			return
		}
		tmpl, err := template.ParseFiles("frontend/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/frontend/static/", controllers.StaticController)
}
