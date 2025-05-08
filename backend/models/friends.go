package models

import (
	"real_time_forum/backend/database"
)

func Friends(userID int) ([]*UserAuth, error) {
	query := `SELECT u.id,u.username, u.firstName, u.lastName,    (
    SELECT MAX(m.sentAt)
    FROM messages m
    WHERE (m.senderID   = ? AND m.receiverID = u.id)
        OR (m.senderID   = u.id  AND m.receiverID = ?)
    ) AS lastAt
        FROM users u
        WHERE u.id != ?
        ORDER BY  lastAt DESC , u.firstName ;
    `
	rows, err := database.DB.Query(query, userID, userID, userID)
	if err != nil {

		return nil, err
	}
	defer rows.Close()

	var users []*UserAuth
	for rows.Next() {
		var user UserAuth
		err := rows.Scan(&user.ID,&user.UserName, &user.FirstName, &user.LastName, &user.LastAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
