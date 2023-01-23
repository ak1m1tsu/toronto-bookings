package models

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

func NewCredentials(body io.Reader) (*Credentials, error) {
	var credentials *Credentials
	if err := json.NewDecoder(body).Decode(&credentials); err != nil {
		return nil, err
	}
	err := ValidateCredentials(credentials)
	return credentials, err
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

func NewApiResponse(status int, body map[string]any) *ApiResponse {
	return &ApiResponse{Status: status, Body: body}
}

func (res *ApiResponse) SetStatus(status int) {
	res.Status = status
}

func (res *ApiResponse) SetBody(body map[string]any) {
	res.Body = body
}

func (res *ApiResponse) SetError(err error) {
	res.Body = map[string]any{"error": err.Error()}
}
