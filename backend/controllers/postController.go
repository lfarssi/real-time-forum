package controllers

import (
	"database/sql"
	"net/http"

	"real_time_forum/backend/utils"
)

func AddPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		utils.ResponseJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"message": "Method not allowed",
			"status":  http.StatusMethodNotAllowed,
		})
	}
	

	// ID := r.Context().Value("userId").(int)
	// title := r.FormValue("title")
	// content := r.FormValue("content")
	// categories := r.Form["category"]
	// file, fileHeader, err := r.FormFile("image")

	// if title == "" || (!requiredContent && content == "") || len(categories) == 0 || len([]rune(content)) > 1000 || len([]rune(title)) > 50 || !models.CheckCatExists(categories, db) {
	// 	handlers.RenderError(w, http.StatusBadRequest)
	// 	return
	// }

	// result, err := db.Exec("INSERT INTO Posts (Title, Content, DateCreation, Image,ID_User) VALUES (?,?,?,?,?)", title, content, time.Now(), image, ID)
	// if err != nil {
	// 	handlers.RenderError(w, http.StatusInternalServerError)
	// 	return
	// }
	// idPost, _ := result.LastInsertId()
	// for _, categoryID := range categories {
	// 	_, err := db.Exec("INSERT INTO PostCategory (ID_Post, ID_Category) VALUES (?, ?)", int(idPost), categoryID)
	// 	if err != nil {
	// 		handlers.RenderError(w, http.StatusInternalServerError)
	// 		return
	// 	}
	// }
	// http.Redirect(w, r, referer, http.StatusFound)
}
