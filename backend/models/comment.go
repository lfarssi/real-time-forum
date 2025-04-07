package models

import (
	"real_time_forum/backend/database"
	"time"
)

type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username"`
	Likes     int       `json:"likes"`
	Dislikes  int       `json:"dislikes"`
	IsLiked 	bool 
	IsDisliked	bool
}

func CreateComment(comment Comment) error {
	query := `
		INSERT INTO comments
		(content , user_id, date_creation, post_id)
		VALUES (?,?,?,?);
	
	`
	_, err := database.Database.Exec(query, comment.Content, comment.UserID, comment.CreatedAt, comment.PostID)
	if err != nil {
		return err
	}
	return nil
}

func GetComments(postid int) ([]Comment, error) {
	query := `
	SELECT c.id, c.content, c.user_id, c.date_creation,  c.post_id , u.username
	FROM comments c
	INNER JOIN  users u 
	ON u.id = c.user_id 
	WHERE post_id = ?
	 ORDER BY date_creation DESC;	
	`
	rows, err := database.Database.Query(query, postid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []Comment
	for rows.Next() {

		var c Comment
		if err := rows.Scan(&c.ID, &c.Content, &c.UserID, &c.CreatedAt, &c.PostID, &c.Username); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}