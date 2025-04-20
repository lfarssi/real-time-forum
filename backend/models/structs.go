package models

type UserAuth struct {
	UserName  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
	Password  string `json:"password"`
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
}


type React struct {
	PostID    int
	CommentID int
	UserID    int
	Status    string
}