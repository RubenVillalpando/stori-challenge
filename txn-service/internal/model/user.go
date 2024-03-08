package model

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type NewUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
