package models

import (
	"time"

	"real_time_forum/backend/database"
)

func Register(username, email, firstName, lastName, gender, password , age string) (int, error) {
	var ID int
	err := database.DB.QueryRow(`INSERT INTO users (Email, UserName, First_Name, Last_Name, Age, Gender, Password, Created_At, Session, Expared_At) 
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING ID`,
		email, username, firstName, lastName, age, gender, password, time.Now().UTC(), "", nil).Scan(&ID)
	if err != nil {
		return -1, err
	}

	return ID, nil
}
