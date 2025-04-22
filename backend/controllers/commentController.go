package controllers

import (
	"encoding/json"
	"net/http"

	"real_time_forum/backend/models"
	"real_time_forum/backend/utils"
)

func GetCommnetsController(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ResponseJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"message": "Method Not Allowed",
			"status":  http.StatusMethodNotAllowed,
		})
		return
	}

	postID := r.URL.Query().Get("postID")

	comments, err := models.GetCommnets(postID)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Error Getting Comments",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	utils.ResponseJSON(w, http.StatusOK, map[string]any{
		"message": "Comments retrieved successfully",
		"status":  http.StatusOK,
		"data":    comments,
	})
}

func AddCommentController(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ResponseJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"message": "Method Not Allowed",
			"status":  http.StatusMethodNotAllowed,
		})
		return
	}

	var comment *models.Comment

	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": err.Error(),
			"status":  http.StatusInternalServerError,
		})
		return
	}

	if comment.Content == "" {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]any{
			"message": "Comment cannot be empty",
			"status":  http.StatusBadRequest,
		})
		return
	}

	comment.UserID = r.Context().Value("userId").(int)

	err := models.AddComment(comment)
	if err != nil {
		if err.Error() == "Post does not exist" {
			utils.ResponseJSON(w, http.StatusUnprocessableEntity, map[string]any{
				"message": err.Error(),
				"status":  http.StatusUnprocessableEntity,
			})
			return
		}

		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": err.Error(),
			"status":  http.StatusInternalServerError,
		})
		return
	}

	utils.ResponseJSON(w, http.StatusOK, map[string]any{
		"message": "Comment added successfully!",
		"status":  http.StatusOK,
	})
}
