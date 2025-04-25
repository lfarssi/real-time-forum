package models

import "real_time_forum/backend/database"

func GetMessage(sender int, receiver int)([]*Message ,error)  {
	query := `SELECT m.id,u.username, m.content , m.sentAt, m.status
	FROM messages m 
	INNER JOIN users u 
	ON u.id=m.receiverID
	WHERE senderID=? AND receiverID=?
	`
	rows, err := database.DB.Query(query, sender, receiver) 
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		var msg Message
		err := rows.Scan(&msg.ID, &msg.Content, &msg.SentAt, &msg.Status, &msg.Username)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil	
}