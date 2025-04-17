package models

import "time"

type GetToken struct {
	Token string `json:"access_token"`
}

type AskLink struct {
	UserID   int
	AuthType string
}

type GithubUser struct {
	UserName string `json:"login"`
	Email    string `json:"email"`
}

type GoogleUser struct {
	UserName string `json:"name"`
	Email    string `json:"email"`
}

type GitHubEmail struct {
	Email   string `json:"email"`
	Primary bool   `json:"primary"`
}

// errors
type ErrorRegister struct {
	ErrName     string
	ErrEmail    string
	ErrPassword string
}

type HomePage struct {
	IsLogged   bool
	UserName   string
	Categories []*Category
	PostCat    []*PostCat
	TotalLikes int
}

type Category struct {
	ID           int
	NameCategory string
}

// allpost
type PostCat struct {
	ID              int
	Title           string
	Content         string
	Image           []byte
	DateCreation    time.Time
	Date            string
	Categories      string
	CreatedBy       string
	Comments        []*Comment
	NemberOfComment int
	Like            int
	Dislike         int
	UserLiked       bool
	UserDisliked    bool
}

type Likes struct {
	ID     int
	IDPost int
}

type Comment struct {
	ID            string
	Content       string
	CreatedBy     string
	DateCreation  time.Time
	Date          string
	Like          int
	Dislike       int
	CUserLiked    bool
	CUserDisliked bool
}
