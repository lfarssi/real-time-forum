package models

import (
	"net/http"
	"real_time_forum/backend/database"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserName             string `json:"username"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	ConfirmationPassword string `json:"confirmationPassword"`
	Role                 string
}

func Login(userName, password string) (int, map[string]string) {
	query := "SELECT id, password FROM users WHERE username = ?"
	statement, err := database.Database.Prepare(query)
	if err != nil {
		return 0, map[string]string{"error": "database error"}
	}
	defer statement.Close()
	var id int
	var hashedPassword string
	err = statement.QueryRow(userName).Scan(&id, &hashedPassword)
	if err != nil {
		return 0, map[string]string{"username": "Username not found"}
	}
	if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) != nil {
		return 0, map[string]string{"password": "Password Incorrect"}
	}
	return id, map[string]string{}
}

func Register(user User) (int, map[string]string) {
	query := "INSERT INTO users (username, email, password) VALUES (?, ?, ?)"
	stm, err :=database.Database.Prepare(query)
	if err != nil {
		return 0, map[string]string{"error": "error preparing query"}
	}
	defer stm.Close()
	res, err := stm.Exec(user.UserName, user.Email, user.Password)
	if err != nil {
		if strings.Contains(err.Error(), "username") {
			return 0, map[string]string{"username": "username already exists"}
		} else if strings.Contains(err.Error(), "email") {
			return 0, map[string]string{"email": "email already exists"}
		}
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, map[string]string{"error": "error getting last id"}
	}
	return int(id), map[string]string{}
}


func GetUserId(r *http.Request) (int, error) {
	var userId int
	token, err := r.Cookie("token")
	if err != nil || token.Value == "" {
		return 0, err
	}
	value := token.Value
	query := "SELECT user_id FROM sessionss WHERE token = ?"
	stm, err := database.Database.Prepare(query)
	if err != nil {
		return 0, err
	}
	err = stm.QueryRow(value).Scan(&userId)
	return int(userId), err
}