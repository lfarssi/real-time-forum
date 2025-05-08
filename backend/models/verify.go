package models

import (
	"real_time_forum/backend/database"
)

// VerifyEmail checks if the given email exists in the database and returns the user ID if found.
func VerifyEmail(Email string) (int, error) {
	var id int
	err := database.DB.QueryRow("SELECT id FROM users WHERE email = ? OR username=?", Email, Email).Scan(&id)
	if err != nil {
		return -1, err
	}
	

	return id, nil
}

// GetPassword retrieves the hashed password for a given user ID from the database.
func GetPassword(id int) (string, error) {
	rows, err := database.DB.Query("SELECT password FROM users WHERE id = ?", id)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	if rows.Next() {
		pass := ""
		err := rows.Scan(&pass)
		if err != nil {
			return "", err
		}
		return pass, nil
	}

	return "", nil
}
