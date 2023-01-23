package models

import "testing"

func TestNewUser(t *testing.T) {
	user, err := NewUser("a@gmail.com", "test1234")
	if err != nil {
		t.Error(err)
	}
	if user.EncryptedPassword == "" {
		t.Errorf("encrypt was corrupted. got: %s", user.EncryptedPassword)
	}
}

func TestUser_ValidatePassword(t *testing.T) {
	pw := "test1234"
	user, err := NewUser("a@gmail.com", pw)
	if err != nil {
		t.Error(err)
	}
	if !user.ValidatePassword(pw) {
		t.Errorf("password validation failed. pw = test1234")
	}
	if user.ValidatePassword("test4321") {
		t.Errorf("password validation failed. pw = test4321")
	}
}
