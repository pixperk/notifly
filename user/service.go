package user

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/pixperk/notifly/user/auth"
	db "github.com/pixperk/notifly/user/db/sqlc"
	"github.com/pixperk/notifly/user/util"
)

type AuthResponse struct {
	Identifier string `json:"identifier"`
	Token      string `json:"token"`
}

type ValidateTokenResponse struct {
	Identifier string    `json:"identifier"`
	UserID     uuid.UUID `json:"user_id"`
}

type Service interface {
	SignUp(ctx context.Context, name, identifier, password string) (*AuthResponse, error)
	SignIn(ctx context.Context, identifier, password string) (*AuthResponse, error)
	ValidateToken(ctx context.Context, token string) (*ValidateTokenResponse, error)
}

type userService struct {
	store      Store
	tokenMaker auth.TokenMaker
	config     *Config
}

func NewService(store Store, config *Config) (Service, error) {

	byteSymmetricKey := []byte(config.TokenSymmetricKey)
	if len(byteSymmetricKey) == 0 {
		return nil, errors.New("token symmetric key is not set")
	}

	tokenMaker, err := auth.NewPasetoMaker(byteSymmetricKey)
	if err != nil {
		return nil, err
	}

	return &userService{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}, nil
}

func (s *userService) SignUp(ctx context.Context, name, identifier, password string) (*AuthResponse, error) {

	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		return nil, err
	}

	args := db.CreateUserParams{
		Name:         name,
		Identifier:   identifier,
		PasswordHash: hashedPassword,
	}
	user, err := s.store.CreateUser(ctx, args)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.tokenMaker.CreateToken(user.ID, user.Identifier, time.Duration(s.config.AccessTokenDuration))
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Identifier: user.Identifier,
		Token:      accessToken,
	}, nil
}

func (s *userService) SignIn(ctx context.Context, identifier, password string) (*AuthResponse, error) {
	user, err := s.store.GetUserByIdentifier(ctx, identifier)
	if err != nil {
		return nil, err
	}

	if err := util.VerifyPassword(password, user.PasswordHash); err != nil {
		return nil, err
	}

	accessToken, err := s.tokenMaker.CreateToken(user.ID, user.Identifier, time.Duration(s.config.AccessTokenDuration))
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Identifier: user.Identifier,
		Token:      accessToken,
	}, nil
}

func (s *userService) ValidateToken(ctx context.Context, token string) (*ValidateTokenResponse, error) {
	payload, err := s.tokenMaker.VerifyToken(token)
	if err != nil {
		return nil, err
	}

	res := &ValidateTokenResponse{
		Identifier: payload.Identifier,
		UserID:     payload.UserId,
	}

	return res, nil

}
