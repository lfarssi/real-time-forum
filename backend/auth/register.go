package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"real_time_forum/backend/controllers"
	"real_time_forum/backend/middleware"
	"real_time_forum/backend/models"
	"real_time_forum/backend/utils"
)

// Register handles the user registration process, validates the inputs, checks for unique username and email,
// hashes the password, inserts the user into the database, generates a session token, and sets it as a secure cookie.
func Register(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/sign-up", http.StatusSeeOther)
		return
	}

	var user *models.UserRegister
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Server error",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	if !verify(w, user.UserName, user.Email, user.FirstName, user.LastName, user.Gender, user.Age, user.Password) {
		return
	}

	isExistsUserName, err := models.UserExists(db, user.UserName, " UserName")
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Server error",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	if isExistsUserName {
		utils.ResponseJSON(w, http.StatusConflict, map[string]any{
			"message": "Username Already taken",
			"status":  http.StatusConflict,
		})
		return
	}

	isExistsEmail, err := models.UserExists(db, user.Email, "Email")
	if err != nil {
		controllers.RenderError(w, http.StatusServiceUnavailable)
		return
	}

	if isExistsEmail {
		utils.ResponseJSON(w, http.StatusConflict, map[string]any{
			"message": "Email Already taken",
			"status":  http.StatusConflict,
		})
		return
	}

	// if !Verify(w, isUniqueUserName, isUniqueEmail, Email, UserName, r.FormValue("Password")) {
	// 	return
	// }
	// password, _ := bcrypt.GenerateFromPassword([]byte(r.FormValue("Password")), 10)

	// result, err := db.Exec("INSERT INTO Users (UserName, Email, Password, Created_At, Session, Expared_At) VALUES ( ?,?,?,?,?,?)", UserName, Email, string(password), time.Now(), "", nil)
	// if err != nil {
	// 	controllers.RenderError(w, http.StatusInternalServerError)
	// 	return
	// }

	// ID, err := result.LastInsertId()
	// if err != nil {
	// 	controllers.RenderError(w, http.StatusInternalServerError)
	// 	return
	// }

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

func verify(w http.ResponseWriter, userName, email, firstName, lastName, gender, age, password string) bool {
	if len([]rune(firstName)) > 30 || len([]rune(lastName)) > 30 {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]any{
			"message": "First name and last name must be less than 30 characters.",
		})
		return false
	}

	if !utils.IsValidName(firstName) || !utils.IsValidName(lastName) {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]any{
			"message": "The first name and last name must contain printable characters and numbers.",
		})
		return false
	}

	if len([]rune(userName)) > 30 {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]any{
			"message": "Username must be less than 30 characters.",
		})
		return false
	}

	if !utils.IsValidName(userName) {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]any{
			"message": "The username must contain printable characters and numbers.",
		})
		return false
	}

	if len([]rune(email)) > 50 {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]any{
			"message": "Email must be less than 50 characters.",
		})
		return false
	}

	if !utils.IsValidEmail(email) {
		utils.ResponseJSON(w, http.StatusUnprocessableEntity, map[string]any{
			"message": "Email must be in the format: john@example.com",
		})
		return false
	}

	_, err := strconv.Atoi(age)
	if err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]any{
			"message": "The age must be a number",
		})
		return false
	}

	if gender != "Male" && gender != "Female" {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]any{
			"message": "The gender must be male or female",
		})
		return false
	}

	if len([]rune(password)) < 8 || len([]rune(password)) > 40 {
		utils.ResponseJSON(w, http.StatusUnprocessableEntity, map[string]any{
			"message": "Password must be greater than 8 characters and less than 40 characters",
		})
		return false
	}

	return true
}

func RegisterPage(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		controllers.RenderError(w, http.StatusMethodNotAllowed)
		return
	}

	_, err := middleware.VerifyCookie(r, db)
	if err == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err = controllers.RenderTemplate(w, "register.html", nil, http.StatusOK)
	if err != nil {
		controllers.RenderError(w, http.StatusInternalServerError)
		return
	}
}
