package controllers

// RenderTemplate renders an HTML template with the provided data and status.
// func RenderTemplate(w http.ResponseWriter, page string, data any, status int) error {
// 	temp, err := template.ParseFiles("./frontend/templates/" + page)
// 	if err != nil {
// 		return err
// 	}

// 	w.WriteHeader(status)
// 	err = temp.Execute(w, data)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// // RenderError renders an error page with the given HTTP status.
// func RenderError(w http.ResponseWriter, status int) {
// 	e := struct {
// 		Type   string
// 		Status int
// 	}{
// 		Type:   http.StatusText(status),
// 		Status: status,
// 	}

// 	err := RenderTemplate(w, "error.html", e, status)
// 	if err != nil {
// 		http.Error(w, "Server error", http.StatusInternalServerError)
// 		return
// 	}
// }
