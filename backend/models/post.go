package models

import (
	"net/http"
	"time"

	"real_time_forum/backend/database"
)

func AddPost(w http.ResponseWriter, title, content string, categories []string, ID int) error {
	var postID int
	err := database.DB.QueryRow("INSERT INTO Posts (Title, Content, DateCreation, ID_User) VALUES ($1, $2, $3, $4) RETURNING ID", title, content, time.Now().UTC(), ID).Scan(&postID)
	if err != nil {
		return err
	}

	for _, categoryID := range categories {
		_, err := database.DB.Exec("INSERT INTO PostCategory (ID_Post, ID_Category) VALUES (?, ?)", postID, categoryID)
		if err != nil {
			return err
		}
	}

	return nil
}
