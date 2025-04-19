package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"real_time_forum/backend/models"
	"real_time_forum/backend/utils"
)

func AddPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		utils.ResponseJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"message": "Method not allowed",
			"status":  http.StatusMethodNotAllowed,
		})
		return
	}

	var post *models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Server error",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	if !verifyData(w, post.Title, post.Content, post.Categories) {
		return
	}

	ID := r.Context().Value("userId").(int)

	fmt.Println(ID)

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

	utils.ResponseJSON(w, http.StatusOK, map[string]any{
		"message": "Post added successfully!",
		"status":  http.StatusOK,
	})
}

func verifyData(w http.ResponseWriter, title, content string, categories []string) bool {
	var message models.ValidationMessagesAddPost
	isValid := true

	if title == "" {
		message.TitleMessage = "Title is required."
		isValid = false
	} else if len([]rune(title)) > 100 {
		message.TitleMessage = "Title must be less than or equal to 100 characters."
		isValid = false
	}

	if content == "" {
		message.ContentMessage = "Content is required."
		isValid = false
	} else if len([]rune(content)) > 1000 {
		message.ContentMessage = "Content must be less than or equal to 1000 characters."
		isValid = false
	}

	if len(categories) == 0 {
		message.CategoryMessage = "At least one category is required."
		isValid = false
	} else if len(categories) > 100 {
		message.CategoryMessage = "You can add up to 100 categories only."
		isValid = false
	} else {
		for _, category := range categories {
			if len([]rune(category)) > 50 {
				message.CategoryMessage = "Each category must be less than or equal to 50 characters."
				isValid = false
				break
			}
		}
	}

	if !isValid {
		utils.ResponseJSON(w, http.StatusBadRequest, message)
		return false
	}

	return true
}
