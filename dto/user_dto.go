package dto

type UserLoginDto struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type UserResponseDto struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
}