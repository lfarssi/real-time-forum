package models

import (
	"real_time_forum/backend/database"
)

func GetMessage(sender int, receiver int) ([]*Message, error) {
	query := `SELECT m.id, m.senderID, m.receiverID, u.username, m.content , m.sentAt, m.status
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
		err := rows.Scan(&msg.ID, &msg.SenderID, &msg.RecipientID, &msg.Username, &msg.Content, &msg.SentAt, &msg.Status)
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

func AddMessage(message *Message) error {
	query := `
		INSERT INTO messages (senderID, receiverID, content, sentAt, status) VALUES ($1, $2, $3, $4, $5) RETURNING id 
	`
	err := database.DB.QueryRow(query, &message.SenderID, &message.RecipientID, &message.Content, &message.SentAt, &message.Status).Scan(&message.ID)
	if err != nil {
		return err
	}
	return nil
}
