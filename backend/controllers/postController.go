package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"real_time_forum/backend/models"
	"real_time_forum/backend/utils"
)
func GetPostController(w http.ResponseWriter, r *http.Request) {
	if r.Method!= http.MethodGet {
		utils.ResponseJSON(w, http.StatusOK, map[string]any{
			"message": "Method Not Allowed",
			"status":  http.StatusMethodNotAllowed,
		})
	}
	posts, err := models.GetPosts()
	if err != nil {
		utils.ResponseJSON(w, http.StatusOK, map[string]any{
			"message": "Error Creation Post",
			"status":  http.StatusOK,
			"data":    posts,
		})
		return
	}

	utils.ResponseJSON(w, http.StatusOK, map[string]any{
		"message": "Posts retrieved successfully",
		"status":  http.StatusOK,
		"data":    posts,
	})
}


func AddPostController(w http.ResponseWriter, r *http.Request) {
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

	post.Title = strings.TrimSpace(post.Title)
	post.Content = strings.TrimSpace(post.Content)

	if !verifyPostData(w, post.Title, post.Content, post.Categories) {
		return
	}

	ID := r.Context().Value("userId").(int)

	err := models.AddPost(w, post.Title, post.Content, post.Categories, ID)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Server error",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	utils.ResponseJSON(w, http.StatusOK, map[string]any{
		"message": "Post added successfully!",
		"status":  http.StatusOK,
	})
}

func verifyPostData(w http.ResponseWriter, title, content string, categories []string) bool {
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
			categoryID, err := strconv.Atoi(category)
			if err != nil {
				message.CategoryMessage = "Each category must be a valid number."
				isValid = false
				break
			}

			if !models.IsExistsCategory(categoryID) {
				message.CategoryMessage = "One or more selected categories do not exist."
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
