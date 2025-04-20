package controllers

import (
	"net/http"

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

	// comments, err := models.GetCommnets()
	// if err != nil {
	// 	utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
	// 		"message": "Error Getting Comments",
	// 		"status":  http.StatusInternalServerError,
	// 	})
	// 	return
	// }

	// utils.ResponseJSON(w, http.StatusOK, map[string]any{
	// 	"message": "Comments retrieved successfully",
	// 	"status":  http.StatusOK,
	// 	"data":    comments,
	// })
}
