package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"

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

	isUniqueUserName, err := models.UserExists(db, user.UserName, " UserName")
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Server error",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	if isUniqueUserName {
		utils.ResponseJSON(w, http.StatusConflict, map[string]any{
			"message": "Username Already taken please chose Another",
			"status": http.StatusConflict,
		})
		return
	}

	// isUniqueEmail, err := models.UserExists(db, Email, " Email ")
	// if err != nil {
	// 	controllers.RenderError(w, http.StatusServiceUnavailable)
	// 	return
	// }

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

func Verify(w http.ResponseWriter, isUniqueUserName, isUniqueEmail bool, Email, UserName, Password string) bool {
	if Email == "" {
		// e.ErrName = "Username cannot be emty"
	}
	if UserName == "" {
		// e.ErrName = "Username cannot be emty"
	}
	if !utils.ValidName(UserName) {
		// e.ErrName = "Username cannot conatains a special charachters like: \"@()-.,;...\" except: _"
	}
	if len([]rune(Password)) < 8 || len([]rune(Password)) > 20 {
		// e.ErrPassword = "Password must be greater than 8 characters and less than 20 characters"
	}
	if !isUniqueUserName {
		// e.ErrName = "Username Already taken please chose Another"
	}
	if !isUniqueEmail {
		// e.ErrEmail = "Email Already taken please chose Another"
	}
	if !utils.IsValidEmail(Email) || len([]rune(Email)) > 200 {
		// e.ErrEmail = "Email must be in the format: example@example.example"
	}
	// if e.ErrEmail != "" || e.ErrName != "" || e.ErrPassword != "" {
	// 	controllers.RenderTemplate(w, "register.html", e, http.StatusConflict)
	// 	return false
	// }
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
