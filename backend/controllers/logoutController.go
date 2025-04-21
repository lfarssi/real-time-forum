package controllers

import (
	"database/sql"
	"net/http"

	"real_time_forum/backend/utils"
)

func Logout(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		utils.ResponseJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"message": "Method not allowed",
			"status":  http.StatusMethodNotAllowed,
		})
		return
	}

	cookie := &http.Cookie{Name: "Token", Value: "", MaxAge: -1, HttpOnly: true}
	cookie2 := &http.Cookie{Name: "UserID", Value: "", MaxAge: -1, HttpOnly: true}

	http.SetCookie(w, cookie)
	http.SetCookie(w, cookie2)

	utils.ResponseJSON(w, http.StatusOK, map[string]any{
		"message": "Logout successful",
		"status":  http.StatusOK,
	})
}
