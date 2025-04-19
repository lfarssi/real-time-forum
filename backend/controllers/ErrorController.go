package controllers

import (
	"net/http"
	"real_time_forum/backend/models"
	"text/template"
)

func ErrorController(w http.ResponseWriter, r *http.Request, StatusCode int, message string) {
	w.WriteHeader(StatusCode)
	errPage := "resources/views/error.html"
	tmp, err := template.ParseFiles(errPage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	errType := models.ErrorType{
		Code:    StatusCode,
		Message: http.StatusText(StatusCode),
	}
	if message != "" {
		errType.Message = message
	}
	if err := tmp.Execute(w, errType); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
