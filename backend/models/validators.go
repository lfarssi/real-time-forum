package models

import "real_time_forum/backend/database"

func IsExistsCategory(categoryID int) bool {
	var id int
	err := database.DB.QueryRow("SELECT id FROM category WHERE id = ?", categoryID).Scan(&id)
	if err != nil {
		return false
	}

	return true
}
