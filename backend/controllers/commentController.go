package controllers

import (
	"fmt"
	"net/http"

	"real_time_forum/backend/models"
	"real_time_forum/backend/utils"
)

func GetCommnetsController(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ResponseJSON(w, http.StatusOK, map[string]any{
			"message": "Method Not Allowed",
			"status":  http.StatusMethodNotAllowed,
		})
		return
	}

	postID := r.URL.Query().Get("postID")
	fmt.Println(postID)
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
