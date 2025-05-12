package models

import (
	"fmt"
	"real_time_forum/backend/database"
	"strings"
	"time"
)

func GetMessage(sender, receiver, lastID int) ([]*Message, error) {
	tx, err := database.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `
        SELECT m.id, m.senderID, m.receiverID, u.username, m.content, m.sentAt, m.status
        FROM messages m
        INNER JOIN users u ON u.id = m.senderID
        WHERE ((m.senderID = ? AND m.receiverID = ?) OR (m.senderID = ? AND m.receiverID = ?)) AND m.id < ?
        ORDER BY m.id DESC
        LIMIT 10
    `
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	rows, err := database.DB.Query(query, sender, receiver, receiver, sender, lastID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		messages []*Message
		ids      []int
	)
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.SenderID, &msg.RecipientID, &msg.Username, &msg.Content, &msg.SentAt, &msg.Status); err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
		ids = append(ids, msg.ID)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(ids) > 0 {
		// Prepare placeholders for the IN clause: "?, ?, ?, ..."
		ph := make([]string, len(ids))
		args := make([]any, len(ids)+1) // +1 for receiverID param

		for i, id := range ids {
			ph[i] = "?"
			args[i] = id
		}

		// Add receiverID as the last argument
		args[len(ids)] = sender

		// Build the query with placeholders and receiverID condition
		upd := fmt.Sprintf(
			"UPDATE messages SET status = 'read' WHERE id IN (%s) AND receiverID = ? AND status = 'unread'",
			strings.Join(ph, ","),
		)

		// Execute the update query with arguments
		if _, err := database.DB.Exec(upd, args...); err != nil {
			return nil, err
		}
		
	}

	return messages, nil
}

func AddMessage(message *Message) error {
	query := `
		INSERT INTO messages (senderID, receiverID, content, sentAt, status) VALUES ($1, $2, $3, $4, $5) RETURNING id 
	`
	err := database.DB.QueryRow(query, &message.SenderID, &message.RecipientID, &message.Content, time.Now(), "unread").Scan(&message.ID)
	if err != nil {
		return err
	}
	return nil
}
func UpdateMessage(senderID, receiverID int, status string) error {
	query := `
        UPDATE messages
        SET status = $1
        WHERE senderID = $2 AND receiverID = $3 AND status = 'unread'
    `
	_, err := database.DB.Exec(query, status, senderID, receiverID)
	return err
}


func GetLastMessageID() (int, error) {
	query := `
		SELECT m.id FROM messages m
		ORDER BY m.id DESC 
		LIMIT 1;
	`

	var id int
	err := database.DB.QueryRow(query).Scan(&id)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return -1, err
	}

	return id + 1, nil
}

func GetUnreadCountsPerFriend(userID, senderId int) (map[int]int, error) {
	query := `
        SELECT senderID, COUNT(*) 
        FROM messages 
        WHERE receiverID = ? AND senderID= ? AND status = 'unread'
        GROUP BY senderID
    `
	rows, err := database.DB.Query(query, userID, senderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[int]int)
	for rows.Next() {
		var friendID, count int
		if err := rows.Scan(&friendID, &count); err != nil {
			return nil, err
		}
		counts[friendID] = count
	}
	return counts, nil
}
func GetUnreadCountsPerFriend2(userID int) (map[int]int, error) {
	query := `
        SELECT senderID, COUNT(*) 
        FROM messages 
        WHERE receiverID = ?  AND status = 'unread'
        GROUP BY senderID
    `
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[int]int)
	for rows.Next() {
		var friendID, count int
		if err := rows.Scan(&friendID, &count); err != nil {
			return nil, err
		}
		counts[friendID] = count
	}
	return counts, nil
}

