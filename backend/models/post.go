package models

import (
	"net/http"
	"time"

	"real_time_forum/backend/database"

)
func GetPosts() ([]*Post, error) {
	query := `
    SELECT p.id, p.user_id, p.title, p.content, GROUP_CONCAT(c.name) AS categories, p.creat_at, u.username
    FROM posts p
    INNER JOIN users u ON p.user_id = u.id
    INNER JOIN post_categorie pc ON p.id = pc.post_id
    INNER JOIN categories c ON pc.categorie_id = c.id
    GROUP BY p.id
    ORDER BY p.creat_at DESC;
    `
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		var post *Post
		var CreatedAt time.Time
		var categorie string
		err = rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content,  &categorie, &CreatedAt, &post.Username)
		if err != nil {
			return nil, err
		}
		post.Categories = append(post.Categories, categorie)
		post.CreatedAt = CreatedAt.Format("2006-01-02 15:04:05")
		posts = append(posts, post)
	}
	return posts, nil
}

func AddPost(w http.ResponseWriter, title, content string, categories []string, ID int) error {
	var postID int
	err := database.DB.QueryRow("INSERT INTO Posts (Title, Content, DateCreation, ID_User) VALUES ($1, $2, $3, $4) RETURNING ID", title, content, time.Now().UTC(), ID).Scan(&postID)
	if err != nil {
		return err
	}

	for _, categoryID := range categories {
		_, err := database.DB.Exec("INSERT INTO PostCategory (ID_Post, ID_Category) VALUES (?, ?)", postID, categoryID)
		if err != nil {
			return err
		}
	}

	return nil
}
