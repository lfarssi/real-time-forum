package models

import (
	"database/sql"
	"time"
)

type UserAuth struct {
	ID        int            `json:"id"`
	UserName  string         `json:"username"`
	Email     string         `json:"email"`
	FirstName string         `json:"firstName"`
	LastName  string         `json:"lastName"`
	Age       int            `json:"age"`
	Gender    string         `json:"gender"`
	Password  string         `json:"password"`
	LastAt    sql.NullString `json:"lastAt"`
	IsOnline  bool           `json:"isOnline"`
}

type ValidationMessagesRegister struct {
	UserNameMessage  string `json:"username,omitempty"`
	EmailMessage     string `json:"email,omitempty"`
	FirstNameMessage string `json:"firstName,omitempty"`
	LastNameMessage  string `json:"lastName,omitempty"`
	GenderMessage    string `json:"gender,omitempty"`
	AgeMessage       string `json:"age,omitempty"`
	PasswordMessage  string `json:"password,omitempty"`
}

type Post struct {
	ID           int      `json:"id"`
	UserID       int      `json:"userID"`
	Title        string   `json:"title"`
	Content      string   `json:"content"`
	Username     string   `json:"username"`
	DateCreation string   `json:"dateCreation"`
	Categories   []string `json:"categories"`
	Likes        int
	Dislikes     int
	IsLiked      bool
	IsDisliked   bool
}

type ValidationMessagesAddPost struct {
	TitleMessage    string `json:"title,omitempty"`
	ContentMessage  string `json:"content,omitempty"`
	CategoryMessage string `json:"category,omitempty"`
}

type Comment struct {
	ID           int    `json:"id"`
	UserID       int    `json:"userID"`
	PostID       int    `json:"postID"`
	Username     string `json:"username"`
	Content      string `json:"content"`
	DateCreation string `json:"dateCreation"`
	Likes        int
	Dislikes     int
	IsLiked      bool
	IsDisliked   bool
}

type React struct {
	PostID    int    `json:"postID"`
	CommentID int    `json:"commentID"`
	UserID    int    `json:"userID"`
	Sender    string `json:"sender"`
	Status    string `json:"status"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Message struct {
	ID          int    `json:"id"`
	SenderID    int    `json:"senderID"`
	RecipientID int    `json:"recipientID"`
	Username    string `json:"username"`
	Content     string `json:"content"`
	SentAt      time.Time   `json:"sentAT"`
	Status      string `json:"status"`
	Type        string `json:"type"`
	LastID      int    `json:"lastID"`
}
