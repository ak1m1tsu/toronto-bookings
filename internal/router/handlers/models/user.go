package models

type UserResponse struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
}

func NewUserResponse(id, email string, isAdmin bool) *UserResponse {
	return &UserResponse{id, email, isAdmin}
}
