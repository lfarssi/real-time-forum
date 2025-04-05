package models

import (
	"real_time_forum/backend/database"
	"time"
)

type Post struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Image         string    `json:"image"`
	Categories    []string  `json:"categories"`
	Likes         int       `json:"likes"`
	Dislikes      int       `json:"dislikes"`
	CreatedAt     string    `json:"created_at"`
	// Comments      []Comment `json:"comments"`
	CommentsCount int       `json:"comments_count"`
	Username      string    `json:"username"`
	IsLiked       bool      `json:"IsLiked"`
	IsDisliked    bool      `json:"IsDisliked"`
}

func GetPost() ([]Post, error)  {
	query := `
    SELECT p.id, p.user_id, p.title, p.content, p.image, GROUP_CONCAT(c.name) AS categories, p.creat_at, u.username
    FROM posts p
    INNER JOIN users u ON p.user_id = u.id
    INNER JOIN post_categorie pc ON p.id = pc.post_id
    INNER JOIN categories c ON pc.categorie_id = c.id
    GROUP BY p.id
    ORDER BY p.creat_at DESC;
    `
	rows, err := database.Database.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var CreatedAt time.Time
		var categorie string
		err = rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.Image, &categorie, &CreatedAt, &post.Username)
		if err != nil {
			return nil, err
		}
		post.Categories = append(post.Categories, categorie)
		post.CreatedAt = CreatedAt.Format("2006-01-02 15:04:05")
		posts = append(posts, post)
	}
	return posts, nil 
}