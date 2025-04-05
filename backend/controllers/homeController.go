package controllers

import (
	"net/http"
	"real_time_forum/backend/models"
	"text/template"
)

func HomePage(w http.ResponseWriter, r *http.Request)  {
	posts , err:= models.GetPost()
	if err!= nil{
		ErrorController(w, r, http.StatusInternalServerError, "Cannot Fetch Post")
		return
	}
	ParseFileController(w, r, "index", posts)
}



func ParseFileController(w http.ResponseWriter, r *http.Request, filename string, data any) {
	filepath := "./frontend/" + filename + ".html"
	components := []string{
		"./frontend/components/header.html",
		"./frontend/components/footer.html",
		"./frontend/components/menu.html",
		"./frontend/components/posts.html",
		"./frontend/components/messages.html",
		// "./frontend/components/register.html",
		// "./frontend/components/login.html",
		"./frontend/components/main.html",
	}

	allFiles := append([]string{filepath}, components...)
	temp, err := template.ParseFiles(allFiles...)
	if err != nil {
		ErrorController(w, r, http.StatusInternalServerError, "Cannot Parse File")
		return
	}
	err1 := temp.Execute(w, data)
	if err1 != nil {
		ErrorController(w, r, http.StatusInternalServerError, "Cannot Execute File")
		return
	}
}