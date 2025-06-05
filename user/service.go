package user

import "context"

type AuthResponse struct {
	Identifier string `json:"identifier"`
	Token      string `json:"token"`
}

type Service interface {
	SignUp(ctx context.Context, name, identifier, password string) (*AuthResponse, error)
	SignIn(ctx context.Context, identifier, password string) (*AuthResponse, error)
}

type userService struct {
	store Store
}
