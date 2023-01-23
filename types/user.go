package types

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID                string `bson:"_id,omitempty"`
	Email             string `bson:"email"`
	EncryptedPassword string `bson:"encryptedPassword"`
	IsAdmin           bool   `bson:"isAdmin"`
}

type CreateUserResponse struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
}

func NewCreateUserResponse(user *User) *CreateUserResponse {
	return &CreateUserResponse{
		ID:      user.ID,
		Email:   user.Email,
		IsAdmin: user.IsAdmin,
	}
}

func NewAdminUser(email, password string) (*User, error) {
	user, err := NewUser(email, password)
	if err != nil {
		return nil, err
	}
	user.IsAdmin = true
	return user, nil
}

func NewUserFromCredentials(creds *Credentials) (*User, error) {
	if err := ValidateCredentials(creds); err != nil {
		return nil, err
	}
	return NewUser(creds.Email, creds.Password)
}

func NewUser(email, password string) (*User, error) {
	epw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &User{
		Email:             email,
		EncryptedPassword: string(epw),
	}, nil
}

func (u *User) ValidatePassword(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(pw)) == nil
}
