package models

import (
	"database/sql"

	"real_time_forum/backend/database"
)

func GetReactionPost(idpost int, status string) ([]*React, error) {
	query := `
		SELECT postID, status, userID FROM postLike
		WHERE postID = ? AND status = ?
	`
	rows, err := database.DB.Query(query, idpost, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reacts []*React
	for rows.Next() {
		var react React
		err := rows.Scan(&react.PostID, &react.Status, &react.UserID)
		if err != nil {
			return nil, err
		}
		reacts = append(reacts, &react)
	}
	return reacts, nil
}

func GetReactionComment(idcomment int, status string) ([]*React, error) {
	query := `
	SELECT commentID, status, userID FROM commentLike
	WHERE commentID = ? AND status = ?
	`
	rows, err := database.DB.Query(query, idcomment, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var reacts []*React
	for rows.Next() {
		var react React
		err := rows.Scan(&react.CommentID, &react.Status, &react.UserID)
		if err != nil {
			continue
		}
		reacts = append(reacts, &react)
	}
	return reacts, nil
}

func InsertReactPost(react React) error {
	react_type, err := ExistReact(react.UserID, react.PostID)
	if err == sql.ErrNoRows {
		query := `
					INSERT INTO postLike (userID, postID, status)
				VALUES(?, ?, ?)
				`
		_, err = database.DB.Exec(query, react.UserID, react.PostID, react.Status)
		if err != nil {
			return err
		}
	} else if err == nil {
		if react_type == react.Status {
			query := `
					DELETE FROM postLike 
					WHERE userID=? AND postID = ?
				`
			_, err := database.DB.Exec(query, react.UserID, react.PostID)
			if err != nil {
				return err
			}
		} else {
			query := `
				UPDATE postLike
				SET status=? 
				WHERE userID=? AND postID=?
			`
			_, err := database.DB.Exec(query, react.Status, react.UserID, react.PostID)
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
		INSERT INTO commentLike (userID, commentID, status)
		VALUES(?, ?, ?)
		`
		_, err = database.DB.Exec(query, react.UserID, react.CommentID, react.Status)
		if err != nil {
			return err
		}
	} else if err == nil {
		if react_type == react.Status {
			query := `
			DELETE FROM commentLike
			WHERE userID=? AND commentID=?
			`
			_, err := database.DB.Exec(query, react.UserID, react.CommentID)
			if err != nil {
				return err
			}
		} else {
			query := `
			UPDATE commentLike
			SET status=? 
			WHERE userID=? AND commentID=?
		`
			_, err := database.DB.Exec(query, react.Status, react.UserID, react.CommentID)
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
		SELECT status 
		FROM postLike
		WHERE userID = ? AND postID= ?
	`
	err := database.DB.QueryRow(query, userId, postId).Scan(&react_type)
	if err != nil {
		return "", err
	}
	return react_type, nil
}

func ExistReactComment(userId, commentId int) (string, error) {
	var react_type string
	query := `
		SELECT status 
		FROM commentLike
		WHERE userID = ? AND commentID= ?
	`
	err := database.DB.QueryRow(query, userId, commentId).Scan(&react_type)
	if err != nil {
		return "", err
	}
	return react_type, nil
}
