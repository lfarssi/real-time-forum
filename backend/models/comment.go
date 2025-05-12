package models

import (
	"time"

	"real_time_forum/backend/database"
)

func GetCommnets(postID string, currentUserID int) ([]*Comment, error) {
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
    SELECT 
        c.id,
        u.username,
        c.content,
        c.dateCreation,
        COALESCE(SUM(CASE WHEN cl.status = 'like' THEN 1 ELSE 0 END), 0) AS likes,
        COALESCE(SUM(CASE WHEN cl.status = 'dislike' THEN 1 ELSE 0 END), 0) AS dislikes,
        MAX(CASE WHEN cl.userID = ? AND cl.status = 'like' THEN 1 ELSE 0 END) = 1 AS isLiked,
        MAX(CASE WHEN cl.userID = ? AND cl.status = 'dislike' THEN 1 ELSE 0 END) = 1 AS isDisliked
    FROM comment c
    INNER JOIN users u ON c.userID = u.id
    LEFT JOIN commentLike cl ON c.id = cl.commentID
    WHERE c.postID = ?
    GROUP BY c.id, u.username, c.content, c.dateCreation
    ORDER BY c.dateCreation DESC;
    `

    rows, err := database.DB.Query(query, currentUserID, currentUserID, postID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var comments []*Comment
    for rows.Next() {
        var c Comment
        var createdAt time.Time

        err := rows.Scan(
            &c.ID,
            &c.Username,
            &c.Content,
            &createdAt,
            &c.Likes,
            &c.Dislikes,
            &c.IsLiked,
            &c.IsDisliked,
        )
        if err != nil {
            return nil, err
        }
        err = tx.Commit()
		if err != nil {
			return nil, err
		}
        c.DateCreation = createdAt.Format(time.RFC3339)
        comments = append(comments, &c)
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
