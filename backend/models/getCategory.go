package models

import (
	"database/sql"
)

func IsExistsCategory(categoryID int, db *sql.DB) bool {
	var id int
	err := db.QueryRow("SELECT ID FROM Category WHERE ID = ?", categoryID).Scan(&id)
	if err != nil {
		return false
	}

	return true
}
