package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"real_time_forum/backend/controllers"
	"real_time_forum/backend/models"
	"real_time_forum/backend/utils"

	"golang.org/x/crypto/bcrypt"
)

// Register handles the user registration process, validates the inputs, checks for unique username and email,
// hashes the password, inserts the user into the database, generates a session token, and sets it as a secure cookie.
func Register(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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
			"message": "Server error",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	if !verifyData(w, user.UserName, user.Email, user.FirstName, user.LastName, user.Gender, user.Age, user.Password) {
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

	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Server error",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	result, err := db.Exec(`INSERT INTO users (UserName, Email, First_Name, Last_Name, Age, Gender, Password, Created_At, Session, Expared_At) 
	VALUES ( ?,?,?,?,?,?,?,?,?,?)`,
		user.UserName, user.Email, user.FirstName, user.LastName, user.Age, user.Gender, password, time.Now().UTC(), "", nil)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Server error",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	ID, err := result.LastInsertId()
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Server error",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	token, err := models.GenerateToken(int(ID), db)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Server error",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	cookie := &http.Cookie{Name: "Token", Value: token, MaxAge: 3600, HttpOnly: true}
	http.SetCookie(w, cookie)

	utils.ResponseJSON(w, http.StatusOK, map[string]any{
		"message": "User registered successfully!",
		"status":  http.StatusOK,
	})
}

func verifyData(w http.ResponseWriter, userName, email, firstName, lastName, gender, age, password string) bool {
	var messages models.ValidationMessages
	hasError := false

	if len([]rune(firstName)) > 30 {
		messages.FirstNameMessage = "First name must be less than 30 characters."
		hasError = true
	} else if !utils.IsValidName(firstName) {
		messages.FirstNameMessage = "First name must contain printable characters and numbers."
		hasError = true
	}

	if len([]rune(lastName)) > 30 {
		messages.LastNameMessage = "Last name must be less than 30 characters."
		hasError = true
	} else if !utils.IsValidName(lastName) {
		messages.LastNameMessage = "Last name must contain printable characters and numbers."
		hasError = true
	}

	if len([]rune(userName)) > 30 {
		messages.UserNameMessage = "Username must be less than 30 characters."
		hasError = true
	} else if !utils.IsValidName(userName) || strings.Contains(userName, " ") {
		messages.UserNameMessage = "Username must contain printable characters and numbers."
		hasError = true
	}

	if len([]rune(email)) > 50 {
		messages.EmailMessage = "Email must be less than 50 characters."
		hasError = true
	} else if !utils.IsValidEmail(email) {
		messages.EmailMessage = "Email must be in the format: john@example.com"
		hasError = true
	}

	if _, err := strconv.Atoi(age); err != nil {
		messages.AgeMessage = "The age must be a number."
		hasError = true
	}

	if gender != "male" && gender != "female" {
		messages.GenderMessage = "The gender must be male or female."
		hasError = true
	}

	passLen := len([]rune(password))
	if passLen < 8 || passLen > 40 {
		messages.PasswordMessage = "Password must be between 8 and 40 characters."
		hasError = true
	}

	if hasError {
		utils.ResponseJSON(w, http.StatusBadRequest, messages)
		return false
	}

	return true
}
