package models

import (
	"time"

	"real_time_forum/backend/database"
)

func GetCommnets(postID int) ([]*Comment, error) {
	query := `
	SELECT u.username, c.id, c.content, c.dateCreation
	FROM comment c
	INNER JOIN users u USING (id)
	WHERE c.postID = ?;
	`

	rows, err := database.DB.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*Comment
	for rows.Next() {
		var comment Comment
		var CreatedAt time.Time
		err = rows.Scan(&comment.Username, &comment.ID, &comment.Content, &comment.DateCreation)
		if err != nil {
			return nil, err
		}

		comment.DateCreation = CreatedAt.Format(time.DateTime)
		comments = append(comments, &comment)
	}

	return comments, nil
}
