package types

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                string `bson:"_id,omitempty" json:"id"`
	Email             string `bson:"email" json:"email"`
	EncryptedPassword string `bson:"encryptedPassword" json:"encrypted_password"`
	IsAdmin           bool   `bson:"isAdmin" json:"is_admin"`
	Token             string `bson:"token" json:"_"`
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
