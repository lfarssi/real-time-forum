package models

import (
	"database/sql"
	"net/http"
	"time"

	"real_time_forum/backend/utils"
)

func AddPost(w http.ResponseWriter, title, content string, categories []string, ID int, db *sql.DB) {
	result, err := db.Exec("INSERT INTO Posts (Title, Content, DateCreation,ID_User) VALUES (?,?,?,?)", title, content, time.Now().UTC(), ID)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Server error",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	postID, _ := result.LastInsertId()
	for _, categoryID := range categories {
		_, err := db.Exec("INSERT INTO PostCategory (ID_Post, ID_Category) VALUES (?, ?)", int(postID), categoryID)
		if err != nil {
			utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
				"message": "Server error",
				"status":  http.StatusInternalServerError,
			})
			return
		}
	}
}
