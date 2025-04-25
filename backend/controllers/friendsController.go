package controllers

import (
	"net/http"

	"real_time_forum/backend/models"
	"real_time_forum/backend/utils"
)

func FriendsController(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		utils.ResponseJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"message": "Method Not Allowed",
			"status":  http.StatusMethodNotAllowed,
		})
		return
	}
	friends, err := models.Friends()
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Cannot Getting Friends",
			"status":  http.StatusInternalServerError,
		})
		return
	}
	utils.ResponseJSON(w, http.StatusOK, map[string]any{
		"message": "Post added successfully!",
		"status":  http.StatusOK,
		"data":    friends,
	})
}
