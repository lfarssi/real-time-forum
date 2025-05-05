package models

import (
	"time"

	"real_time_forum/backend/database"
)

func GetCommnets(postID string) ([]*Comment, error) {
	query := `
	SELECT u.username, c.id, c.content, c.dateCreation
	FROM comment c
	INNER JOIN users u ON c.userID = u.ID
	WHERE c.postID = ?
	ORDER BY c.dateCreation desc;

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

        // Get like and dislike counts
        likes, err := GetReactionComment(comment.ID, "like")
        if err != nil {
            return nil, err
        }
        dislikes, err := GetReactionComment(comment.ID, "dislike")
        if err != nil {
            return nil, err
        }

        comment.Likes = len(likes)
        comment.Dislikes = len(dislikes)
		comments = append(comments, &comment)
	}

	return comments, nil
}

func AddComment(Comment *Comment) error {
	query := `
		INSERT INTO comment
		(content , userID, dateCreation, postID)
		VALUES (?,?,?,?);
	
	`
	_, err := database.DB.Exec(query, Comment.Content, Comment.UserID, time.Now().UTC(), Comment.PostID)
	if err != nil {
		return err
	}
	return nil
}
