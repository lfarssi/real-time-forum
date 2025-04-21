package controllers

import (
	"net/http"
	"real_time_forum/backend/models"
	"real_time_forum/backend/utils"
)

func CategoryController(w http.ResponseWriter, r *http.Request) {
	categories, err := models.GetCategories()
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": err.Error(),
			"status":  http.StatusInternalServerError,
		})
		return
	}
	utils.ResponseJSON(w, http.StatusOK, map[string]any{
		"message": "Category retrieved successfully",
		"status":  http.StatusOK,
		"data": categories,
	})

}
