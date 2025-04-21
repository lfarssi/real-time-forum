package router

import (
	"html/template"
	"net/http"

	"real_time_forum/backend/controllers"
)

func WebRouter() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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

	http.HandleFunc("/getCategory", controllers.CategoryController)
	http.HandleFunc("/frontend/static/", controllers.StaticController)
}
