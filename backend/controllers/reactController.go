package controllers

import (
	"encoding/json"
	"net/http"

	"real_time_forum/backend/models"
	"real_time_forum/backend/utils"
)

func ReactPostController(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		utils.ResponseJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"message": "Method Not Allowed",
			"status":  http.StatusMethodNotAllowed,
		})
		return
	}

	var react models.React
	var err error

	if err = json.NewDecoder(r.Body).Decode(&react); err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Server error",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	react.UserID = r.Context().Value("userId").(int)

	if react.Sender == "post" {
		err = models.InsertReactPost(react)
		if err != nil {
			utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
				"message": "Cannot react Try Again",
				"status":  http.StatusInternalServerError,
			})
			return
		}

	} else if react.Sender == "comment" {
		err = models.InsertReactComment(react)
		if err != nil {
			utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
				"message": "Cannot react Try Again",
				"status":  http.StatusInternalServerError,
			})
			return
		}

	}

	utils.ResponseJSON(w, http.StatusOK, map[string]any{
		"message": "React added successfully!",
		"status":  http.StatusOK,
	})
}
