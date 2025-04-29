package controllers

import (
	"net/http"

	"real_time_forum/backend/utils"
)

func LogoutController(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ResponseJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"message": "Method not allowed",
			"status":  http.StatusMethodNotAllowed,
		})
		return
	}

	cookie := &http.Cookie{Name: "Token", Value: "", MaxAge: -1, Path: "/"}
	// cookie2 := &http.Cookie{Name: "UserID", Value: "", MaxAge: -1, HttpOnly: true}

	http.SetCookie(w, cookie)
	// http.SetCookie(w, cookie2)

	utils.ResponseJSON(w, http.StatusOK, map[string]any{
		"message": "Logout successful",
		"status":  http.StatusOK,
	})
}
