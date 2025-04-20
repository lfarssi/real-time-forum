package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"real_time_forum/backend/models"
	"real_time_forum/backend/utils"

	"golang.org/x/crypto/bcrypt"
)

// Register handles the user registration process, validates the inputs, checks for unique username and email,
// hashes the password, inserts the user into the database, generates a session token, and sets it as a secure cookie.
func RegisterController(w http.ResponseWriter, r *http.Request) {
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

	if !verifyRegisterData(w, user.UserName, user.Email, user.FirstName, user.LastName, user.Gender, user.Password) {
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

	ID, err := models.Register(user.UserName, user.Email, user.FirstName, user.LastName, user.Gender, string(password), user.Age)
	if err != nil {
		if strings.Contains(err.Error(), "UserName") {
			utils.ResponseJSON(w, http.StatusConflict, map[string]any{
				"message": "Username already exists",
				"status":  http.StatusConflict,
			})
			return
		} else if strings.Contains(err.Error(), "Email") {
			utils.ResponseJSON(w, http.StatusConflict, map[string]any{
				"message": "Email already exists",
				"status":  http.StatusConflict,
			})
			return
		}

		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Server error",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	token, err := models.GenerateToken(int(ID))
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

func verifyRegisterData(w http.ResponseWriter, userName, email, firstName, lastName, gender, password string) bool {
	var messages models.ValidationMessagesRegister
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

	if !utils.IsValidUserName(userName) {
		messages.UserNameMessage = "Username must contain printable characters and numbers and between 3 and 13 character."
		hasError = true
	}

	if len([]rune(email)) > 50 {
		messages.EmailMessage = "Email must be less than 50 characters."
		hasError = true
	} else if !utils.IsValidEmail(email) {
		messages.EmailMessage = "Email must be in the format: john@example.com"
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
