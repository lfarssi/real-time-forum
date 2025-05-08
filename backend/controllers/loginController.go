package controllers

import (
	"encoding/json"
	"net/http"

	"real_time_forum/backend/models"
	"real_time_forum/backend/utils"

	"golang.org/x/crypto/bcrypt"
)

// LoginPage renders the login page, or redirects if the user is already authenticated.
func LoginController(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ResponseJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"message": "Method not allowed",
			"status":  http.StatusMethodNotAllowed,
		})
		return
	}

	var user *models.UserAuth
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": err.Error(),
			"status":  http.StatusInternalServerError,
		})
		return
	}

	ID, err := models.VerifyEmail(user.Email)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": err.Error(),
			"status":  http.StatusInternalServerError,
		})
		return
	}

	if ID == 0 {
		utils.ResponseJSON(w, http.StatusUnprocessableEntity, map[string]any{
			"message": "Incorrect email or password",
			"status":  http.StatusUnprocessableEntity,
		})
		return
	}

	password, err := models.GetPassword(int(ID))
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": err.Error(),
			"status":  http.StatusInternalServerError,
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password))
	if err != nil {
		utils.ResponseJSON(w, http.StatusUnprocessableEntity, map[string]any{
			"message": "Incorrect email or password",
			"status":  http.StatusUnprocessableEntity,
		})
		return
	}

	token, err := models.GenerateToken(int(ID))
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": err.Error(),
			"status":  http.StatusInternalServerError,
		})
		return
	}

	cookie := &http.Cookie{Name: "Token", Value: token, MaxAge: 3600 * 24, HttpOnly: true, Path: "/"}
	http.SetCookie(w, cookie)

	utils.ResponseJSON(w, http.StatusOK, map[string]any{
		"message": "User logged in successfully!",
		"status":  http.StatusOK,
	})
}
