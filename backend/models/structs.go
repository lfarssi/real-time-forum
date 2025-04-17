package models

type UserRegister struct {
	UserName  string `json:"userName"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       string `json:"age"`
	Gender    string `json:"gender"`
	Password  string `json:"password"`
}

type ValidationMessages struct {
	UserNameMessage  string `json:"username,omitempty"`
	EmailMessage     string `json:"email,omitempty"`
	FirstNameMessage string `json:"firstName,omitempty"`
	LastNameMessage  string `json:"lastName,omitempty"`
	GenderMessage    string `json:"gender,omitempty"`
	AgeMessage       string `json:"age,omitempty"`
	PasswordMessage  string `json:"password,omitempty"`
}
