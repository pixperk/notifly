package auth

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalid")
)

type Payload struct {
	ID         uuid.UUID `json:"id"`
	UserId     uuid.UUID `json:"user_id"`
	Identifier string    `json:"identifier"`
	IssuedAt   time.Time `json:"issued_at"`
	ExpiredAt  time.Time `json:"expired_at"`
}

func NewPayload(userId uuid.UUID, identifier string, duration time.Duration) (*Payload, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &Payload{
		ID:         tokenId,
		UserId:     userId,
		Identifier: identifier,
		IssuedAt:   time.Now(),
		ExpiredAt:  time.Now().Add(duration),
	}, nil
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return ErrExpiredToken
	}

	if p.ID == uuid.Nil || p.UserId == uuid.Nil {
		return ErrInvalidToken
	}

	return nil
}
