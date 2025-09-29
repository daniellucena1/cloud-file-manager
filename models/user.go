package models

type User struct {
	ID int `json:"id"`
	Name string `json:"name_user"`
	Email string `json:"email_user"`
	Password string `json:"password_user"`
}

