package service

import (
	"errors"
	"testing"

	"crud_alumni/app/model"
	"crud_alumni/app/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// helper: buat hash bcrypt
func hashPassword(t *testing.T, plain string) string {
	t.Helper()
	h, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	return string(h)
}

func TestLogin_Success(t *testing.T) {
	// Mock repository function supaya tidak mengakses MongoDB nyata
	repository.FindUserByUsernameOrEmailFunc = func(identifier string) (*model.User, string, error) {
		u := &model.User{
			ID:        primitive.NewObjectID(),
			Username:  "alice",
			Email:     "alice@example.com",
			Role:      "user",
			CreatedAt: "2025-01-01",
		}
		hashed := hashPassword(t, "supersecret")
		return u, hashed, nil
	}
	// restore ke nil setelah test supaya tidak mempengaruhi test lain
	defer func() { repository.FindUserByUsernameOrEmailFunc = nil }()

	req := model.LoginRequest{
		Username: "alice",
		Password: "supersecret",
	}

	resp, err := Login(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == nil {
		t.Fatalf("expected response, got nil")
	}
	if resp.User.Username != "alice" {
		t.Errorf("expected username alice, got %s", resp.User.Username)
	}
	if resp.Token == "" {
		t.Errorf("expected non-empty token")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	repository.FindUserByUsernameOrEmailFunc = func(identifier string) (*model.User, string, error) {
		u := &model.User{
			ID:       primitive.NewObjectID(),
			Username: "bob",
			Email:    "bob@example.com",
			Role:     "user",
		}
		hashed := hashPassword(t, "correctpassword")
		return u, hashed, nil
	}
	defer func() { repository.FindUserByUsernameOrEmailFunc = nil }()

	req := model.LoginRequest{
		Username: "bob",
		Password: "wrongpassword",
	}
	_, err := Login(req)
	if err == nil {
		t.Fatalf("expected error for wrong password, got nil")
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	repository.FindUserByUsernameOrEmailFunc = func(identifier string) (*model.User, string, error) {
		return nil, "", errors.New("not found")
	}
	defer func() { repository.FindUserByUsernameOrEmailFunc = nil }()

	req := model.LoginRequest{
		Username: "nonexistent",
		Password: "whatever",
	}
	_, err := Login(req)
	if err == nil {
		t.Fatalf("expected error when user not found, got nil")
	}
}
