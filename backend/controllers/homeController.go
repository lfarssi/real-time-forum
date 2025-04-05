package controllers


import (
	"net/http"
	"text/template"
)

func HomePage()  {
	
}



func ParseFileController(w http.ResponseWriter, r *http.Request, filename string, data any) {
	filepath := "./resources/views/" + filename + ".html"
	components := []string{
		"./resources/views/components/navbar.html",
		"./resources/views/components/footer.html",
		"./resources/views/components/menu.html",
		"./resources/views/components/posts.html",
		"./resources/views/components/displayPost.html",
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