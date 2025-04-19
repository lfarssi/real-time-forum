package models

type UserAuth struct {
	UserName  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       string `json:"age"`
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
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Categories []string `json:"category"`
}

type ValidationMessagesAddPost struct {
	TitleMessage    string `json:"title,omitempty"`
	ContentMessage  string `json:"content,omitempty"`
	CategoryMessage string `json:"category,omitempty"`
}
