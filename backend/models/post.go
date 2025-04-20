package models

import (
	"net/http"
	"time"

	"real_time_forum/backend/database"
)

func GetPosts() ([]*Post, error) {
	query := `
    SELECT p.id, p.userID, p.title, p.content, GROUP_CONCAT(c.name) AS categories, p.dateCreation, u.username
    FROM posts p
    INNER JOIN users u ON p.userID = u.id
    INNER JOIN postCategory pc ON p.id = pc.postID
    INNER JOIN category c ON pc.categoryID = c.id
    GROUP BY p.id
    ORDER BY p.dateCreation DESC;
    `
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		var post Post
		var CreatedAt time.Time
		var categorie string
		err = rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &categorie, &CreatedAt, &post.Username)
		if err != nil {
			return nil, err
		}

		post.Categories = append(post.Categories, categorie)
		post.DateCreation = CreatedAt.Format(time.DateTime)
		posts = append(posts, &post)
	}

	return posts, nil
}

func AddPost(w http.ResponseWriter, title, content string, categories []string, ID int) error {
	var postID int
	err := database.DB.QueryRow("INSERT INTO posts (title, content, dateCreation, userID) VALUES ($1, $2, $3, $4) RETURNING id", title, content, time.Now().UTC(), ID).Scan(&postID)
	if err != nil {
		return err
	}

	for _, categoryID := range categories {
		_, err := database.DB.Exec("INSERT INTO postCategory (postID, categoryID) VALUES (?, ?)", postID, categoryID)
		if err != nil {
			return err
		}
	}

	return nil
}
