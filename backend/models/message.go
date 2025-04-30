package models

import (
	"time"

	"real_time_forum/backend/database"
)

func GetMessage(sender int, receiver int) ([]*Message, error) {
	query := `SELECT m.id, m.senderID, m.receiverID, u.username, m.content , m.sentAt, m.status
	FROM messages m 
	INNER JOIN users u 
	ON u.id=m.senderID
	WHERE (senderID=? OR senderID = ?) AND (receiverID=? OR receiverID=?)
	ORDER BY m.id
	LIMIT 10
	`
	rows, err := database.DB.Query(query, sender, receiver, receiver, sender)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		var msg Message
		var t time.Time
		err := rows.Scan(&msg.ID, &msg.SenderID, &msg.RecipientID, &msg.Username, &msg.Content, &t, &msg.Status)
		if err != nil {
			return nil, err
		}
		msg.SentAt = t.Format(time.TimeOnly)
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
	err := database.DB.QueryRow(query, &message.SenderID, &message.RecipientID, &message.Content, time.Now(), &message.Status).Scan(&message.ID)
	if err != nil {
		return err
	}
	return nil
}
