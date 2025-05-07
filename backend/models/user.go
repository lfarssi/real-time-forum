package models

import (
	"database/sql"
	"fmt"
	"time"

	"real_time_forum/backend/database"
)

func Register(username, email, firstName, lastName, gender, password string, age int) (int, error) {
	var ID int
	err := database.DB.QueryRow(`INSERT INTO users (email, username, firstName, lastName, age, gender, password, createdAt, session, expiredAt) 
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING ID`,
		email, username, firstName, lastName, age, gender, password, time.Now().UTC(), "", nil).Scan(&ID)
	if err != nil {
		return -1, err
	}

	return ID, nil
}

func GetUserByID(userID int) (UserAuth, error) {
	var user UserAuth

	query := `SELECT id, firstName, lastName, userName FROM users WHERE id = ?`

	row := database.DB.QueryRow(query, userID)
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.UserName)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("user with id %d not found", userID)
		}
		return user, err
	}

	return user, nil
}
