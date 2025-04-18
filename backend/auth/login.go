package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"real_time_forum/backend/controllers"
	"real_time_forum/backend/middleware"
	"real_time_forum/backend/models"
	"real_time_forum/backend/utils"

	"golang.org/x/crypto/bcrypt"
)

// LoginPage renders the login page, or redirects if the user is already authenticated.
func Login(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return
	}

	var user *models.UserAuth
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Server error",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	ID, _, err := models.VerifyEmail(db, user.Email)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Server error",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	if ID == -1 {
		utils.ResponseJSON(w, http.StatusUnprocessableEntity, map[string]any{
			"message": "Incorrect email or password",
			"status":  http.StatusUnprocessableEntity,
		})
		return
	}

	password, err := models.GetPassword(db, int(ID))
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Server error",
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

	// token, err := models.GenerateToken(int(ID), db)
	// if err != nil {
	// 	controllers.RenderError(w, http.StatusInternalServerError)
	// 	return
	// }

	// cookie := &http.Cookie{Name: "Token", Value: token, MaxAge: 3600, HttpOnly: true}

	// http.SetCookie(w, cookie)
	// http.Redirect(w, r, "/", http.StatusSeeOther)

	utils.ResponseJSON(w, http.StatusOK, map[string]any{
		"message": "User registered successfully!",
		"status":  http.StatusOK,
	})
}

func LoginPage(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		controllers.RenderError(w, http.StatusMethodNotAllowed)
		return
	}

	_, err := middleware.VerifyCookie(r, db)
	if err == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err = controllers.RenderTemplate(w, "login.html", nil, http.StatusOK)
	if err != nil {
		controllers.RenderError(w, http.StatusInternalServerError)
		return
	}
}
