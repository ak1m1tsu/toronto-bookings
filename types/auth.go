package types

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewCredentialsFromRequestBody(body io.Reader) (*Credentials, error) {
	var credentials *Credentials
	if err := json.NewDecoder(body).Decode(&credentials); err != nil {
		return nil, err
	}
	return credentials, nil
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func NewClaims(userID string) *Claims {
	return &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		},
	}
}

type ApiResponse struct {
	Status int            `json:"status"`
	Body   map[string]any `json:"body"`
}

func ValidateCredentials(creds *Credentials) error {
	if !isFieldValid(creds.Email, emailRegexPattern) {
		return fmt.Errorf("email is not valid")
	}
	if len(creds.Password) < 8 {
		return fmt.Errorf("password is not valid")
	}
	return nil
}
