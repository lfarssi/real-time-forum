package models

import (
	"database/sql"
	"real_time_forum/backend/database"
)

type React struct {
	PostID    int
	CommentID int
	UserID    int
	Status    string
	Sender    string
}

func GetReactionPost(idpost int, status string) ([]React, error) {
	var reacts []React
	query := `
		SELECT post_id, react_type, user_id FROM reactPost
		WHERE post_id = ? AND react_type = ?
	`
	rows, err := database.Database.Query(query, idpost, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var react React
		err := rows.Scan(&react.PostID, &react.Status, &react.UserID)
		if err != nil {
			return nil , err
		}
		reacts = append(reacts, react)
	}
	return reacts, nil

}
func GetReactionComment(idcomment int, status string) ([]React, error) {
	var reacts []React
	query := `
		SELECT comment_id, react_type, user_id FROM reactComment
		WHERE comment_id = ? AND react_type = ?
	`
	rows, err := database.Database.Query(query, idcomment, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var react React
		err := rows.Scan(&react.CommentID, &react.Status, &react.UserID)
		if err != nil {
			continue
		}
		reacts = append(reacts, react)
	}
	return reacts, nil

}
func InsertReactPost(react React) error {
	react_type, err := ExistReact(react.UserID, react.PostID)
	if err == sql.ErrNoRows {
		query := `
					INSERT INTO reactPost (user_id, post_id, react_type)
				VALUES(?, ?, ?)
				`
		_, err = database.Database.Exec(query, react.UserID, react.PostID, react.Status)
		if err != nil {
			return err
		}
	} else if err == nil {
		if react_type == react.Status {
			query := `
					DELETE FROM reactPost 
					WHERE user_id=? AND post_id = ?
				`
			_, err := database.Database.Exec(query, react.UserID, react.PostID)
			if err != nil {
				return err
			}
		} else {
			query := `
				UPDATE reactPost
				SET react_type=? 
				WHERE user_id=? AND post_id=?
			`
			_, err := database.Database.Exec(query, react.Status, react.UserID, react.PostID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func InsertReactComment(react React) error {
	react_type, err := ExistReactComment(react.UserID, react.CommentID)
	if err == sql.ErrNoRows {
		query := `
		INSERT INTO reactComment (user_id, comment_id, react_type)
		VALUES(?, ?, ?)
		`
		_, err = database.Database.Exec(query, react.UserID, react.CommentID, react.Status)
		if err != nil {
			return err
		}
	} else if err == nil {
		if react_type == react.Status {
			query := `
			DELETE FROM reactComment
			WHERE user_id=? AND comment_id=?
			`
			_, err := database.Database.Exec(query, react.UserID, react.CommentID)
			if err != nil {
				return err
			}
		} else {
			query := `
			UPDATE reactComment
			SET react_type=? 
			WHERE user_id=? AND comment_id=?
		`
			_, err := database.Database.Exec(query, react.Status, react.UserID, react.CommentID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func ExistReact(userId, postId int) (string, error) {
	var react_type string
	query := `
		SELECT react_type 
		FROM reactPost
		WHERE user_id = ? AND post_id= ?
	`
	err := database.Database.QueryRow(query, userId, postId).Scan(&react_type)
	if err != nil {
		return "", err
	}
	return react_type, nil
}
func ExistReactComment(userId, commentId int) (string, error) {
	var react_type string
	query := `
		SELECT react_type 
		FROM reactComment
		WHERE user_id = ? AND comment_id= ?
	`
	err := database.Database.QueryRow(query, userId, commentId).Scan(&react_type)
	if err != nil {
		return "", err
	}
	return react_type, nil
}