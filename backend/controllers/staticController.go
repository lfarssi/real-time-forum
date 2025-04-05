package controllers
import (
	"net/http"
	"os"
	"strings"
)

func StaticController(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		ErrorController(w, r, http.StatusMethodNotAllowed, "")
		return
	}
	filePath := strings.TrimPrefix(r.URL.Path, "/frontend/")
	fullPath := "frontend/" + filePath

	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			ErrorController(w, r, http.StatusNotFound, "Path Doesn't Exist")
		} else {
			ErrorController(w, r, http.StatusInternalServerError, "")
		}
		return
	}
	if info.IsDir() {
		ErrorController(w, r, http.StatusForbidden, "You can't access this directory")
		return
	}

	fs := http.Dir("frontend")
	http.StripPrefix("/frontend/", http.FileServer(fs)).ServeHTTP(w, r)

}