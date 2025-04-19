package models

import (
	"real_time_forum/backend/database"
)

// UserExists checks if a user exists in the database based on the given value and search criteria.
func UserExists(value string, searchBy string) (bool, error) {
	rows, err := database.DB.Query("SELECT ID from users where "+searchBy+" =?", value)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

// VerifyEmail checks if the given email exists in the database and returns the user ID if found.
func VerifyEmail(Email string) (int64, int, error) {
	rows, err := database.DB.Query("SELECT ID, AuthType FROM users WHERE Email = ?", Email)
	if err != nil {
		return -1, 0, err
	}
	defer rows.Close()

	if rows.Next() {
		var id int64
		var AuthType int
		err := rows.Scan(&id, &AuthType)
		if err != nil {
			return -1, 0, err
		}
		return id, AuthType, nil
	}

	return -1, 0, nil
}

// GetPassword retrieves the hashed password for a given user ID from the database.
func GetPassword(id int) (string, error) {
	rows, err := database.DB.Query("SELECT Password FROM users WHERE ID = ?", id)
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
