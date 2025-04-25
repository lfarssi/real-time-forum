package models

import (
	"database/sql"
	"time"

	"real_time_forum/backend/database"
)

func GetPosts(userID int) ([]*Post, error) {
	query := `
    SELECT p.id, p.userID, p.title, p.content, GROUP_CONCAT(DISTINCT c.name) AS categories, p.dateCreation, u.username
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
		query2 := `SELECT 
		COUNT(CASE WHEN pl.status = 'like' THEN 1 END) AS likeCount,
		COUNT(CASE WHEN pl.status = 'dislike' THEN 1 END) AS dislikeCount,
        MAX(CASE WHEN pl.userID = ? AND pl.status = 'like' THEN 1 ELSE 0 END) AS isLiked,
        MAX(CASE WHEN pl.userID = ? AND pl.status = 'dislike' THEN 1 ELSE 0 END) AS isDisliked
		FROM postLike pl
		WHERE pl.postID = ?
        GROUP BY pl.postID;
	`
		row := database.DB.QueryRow(query2, userID, userID, post.ID)
		err = row.Scan(&post.Likes, &post.Dislikes, &post.IsLiked, &post.IsDisliked)
		if err != nil {
			if err == sql.ErrNoRows {
				post.Likes = 0
				post.Dislikes = 0
				post.IsLiked = false
				post.IsDisliked = false
			} else {
				return nil, err
			}
		}
		post.Categories = append(post.Categories, categorie)
		post.DateCreation = CreatedAt.Format(time.DateTime)
		posts = append(posts, &post)
	}
	return posts, nil
}

func AddPost(post *Post) error {
	var postID int
	err := database.DB.QueryRow("INSERT INTO posts (title, content, dateCreation, userID) VALUES ($1, $2, $3, $4) RETURNING id", post.Title, post.Content, time.Now().UTC(), post.UserID).Scan(&postID)
	if err != nil {
		return err
	}

	for _, categoryID := range post.Categories {
		_, err := database.DB.Exec("INSERT INTO postCategory (postID, categoryID) VALUES (?, ?)", postID, categoryID)
		if err != nil {
			return err
		}
	}

	return nil
}

func LikedPost(userID int) ([]*Post, error) {
	query := `
	SELECT p.id , p.title,p.content,p.dateCreation ,u.username , GROUP_CONCAT(c.name) AS categories
	FROM posts p 
	INNER JOIN users u ON u.id=p.userID
	INNER JOIN postLike r ON p.id = r.postID
	INNER JOIN postCategory pc ON p.id = pc.postID
    INNER JOIN category c ON pc.categoryID = c.id
	WHERE status='like' AND r.userID=?
	GROUP BY p.id
	ORDER BY p.dateCreation DESC
	`
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var LikedPost []*Post
	for rows.Next() {
		var post Post
		var categorie string
		var CreatedAt time.Time
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &CreatedAt, &post.Username, &categorie)
		if err != nil {
			return nil, err
		}
		query2 := `SELECT 
		COUNT(CASE WHEN pl.status = 'like' THEN 1 END) AS likeCount,
		COUNT(CASE WHEN pl.status = 'dislike' THEN 1 END) AS dislikeCount,
        MAX(CASE WHEN pl.userID = ? AND pl.status = 'like' THEN 1 ELSE 0 END) AS isLiked,
        MAX(CASE WHEN pl.userID = ? AND pl.status = 'dislike' THEN 1 ELSE 0 END) AS isDisliked
		FROM postLike pl
		WHERE pl.postID = ?
        GROUP BY pl.postID;
	`
		row := database.DB.QueryRow(query2, userID, userID, post.ID)
		err = row.Scan(&post.Likes, &post.Dislikes, &post.IsLiked, &post.IsDisliked)
		if err != nil {
			if err == sql.ErrNoRows {
				post.Likes = 0
				post.Dislikes = 0
				post.IsLiked = false
				post.IsDisliked = false
			} else {
				return nil, err
			}
		}
		post.Categories = append(post.Categories, categorie)
		post.DateCreation = CreatedAt.Format(time.DateTime)
		LikedPost = append(LikedPost, &post)
	}
	return LikedPost, nil

}

func CreatedPost(userID int) ([]Post, error) {
	query := `
	SELECT p.id , p.title,p.content,p.dateCreation ,u.username , GROUP_CONCAT(DISTINCT c.name) AS categories
	FROM posts p 
	INNER JOIN users u ON u.id=p.userID
	INNER JOIN postCategory pc ON p.id = pc.postID
    INNER JOIN category c ON pc.categoryID = c.id
	WHERE p.userID=?
	GROUP BY p.id
	ORDER BY p.dateCreation DESC
	`
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var createdPost []Post
	for rows.Next() {
		var post Post
		var categorie string
		var CreatedAt time.Time
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &CreatedAt, &post.Username, &categorie)
		if err != nil {
			return nil, err
		}
		query2 := `SELECT 
		COUNT(CASE WHEN pl.status = 'like' THEN 1 END) AS likeCount,
		COUNT(CASE WHEN pl.status = 'dislike' THEN 1 END) AS dislikeCount,
        MAX(CASE WHEN pl.userID = ? AND pl.status = 'like' THEN 1 ELSE 0 END) AS isLiked,
        MAX(CASE WHEN pl.userID = ? AND pl.status = 'dislike' THEN 1 ELSE 0 END) AS isDisliked
		FROM postLike pl
		WHERE pl.postID = ?
        GROUP BY pl.postID;
	`
		row := database.DB.QueryRow(query2, userID, userID, post.ID)
		err = row.Scan(&post.Likes, &post.Dislikes, &post.IsLiked, &post.IsDisliked)
		if err != nil {
			if err == sql.ErrNoRows {
				post.Likes = 0
				post.Dislikes = 0
				post.IsLiked = false
				post.IsDisliked = false
			} else {
				return nil, err
			}
		}
		post.Categories = append(post.Categories, categorie)
		post.DateCreation = CreatedAt.Format(time.DateTime)
		createdPost = append(createdPost, post)
	}
	return createdPost, nil
}

func GetPostsByCategory(idCategory int) ([]Post, error) {
	query := `
	SELECT   p.id, p.title, p.content, c.name, p.dateCreation, u.username
	FROM posts p
	INNER JOIN users u ON p.userID = u.id
	INNER JOIN postCategory pc ON p.id = pc.postID
	INNER JOIN category c ON pc.categoryID = c.id
	WHERE pc.categoryID =?
	ORDER BY p.dateCreation DESC;
	`
	rows, err := database.DB.Query(query, idCategory)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []Post
	tempPosts := make(map[int]*Post)
	for rows.Next() {
		var post Post
		var CreatedAt time.Time
		var categorie string
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &categorie, &CreatedAt, &post.Username)
		if err != nil {
			return nil, err
		}
		if temposts, ok := tempPosts[post.ID]; ok {
			temposts.Categories = append(temposts.Categories, categorie)
		} else {
			post.Categories = CorrectCategories(post.ID)
			tempPosts[post.ID] = &post
		}
		post.DateCreation = CreatedAt.Format(time.DateTime)

	}

	for _, post := range tempPosts {
		posts = append(posts, *post)
	}
	return posts, nil

}

func CorrectCategories(id int) []string {
	query := `SELECT c.name FROM category c
	INNER JOIN postCategory pc ON c.id = pc.categoryID
	WHERE pc.postID = ?`
	rows, err := database.DB.Query(query, id)
	if err != nil {
		return nil
	}
	defer rows.Close()
	categories := []string{}
	for rows.Next() {
		var categorie string
		err := rows.Scan(&categorie)
		if err != nil {
			continue
		}
		categories = append(categories, categorie)
	}
	return categories
}
