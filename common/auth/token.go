package auth

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/google/uuid"
	"github.com/o1egl/paseto"
)

type TokenMaker interface {
	CreateToken(userId uuid.UUID, identifier string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey []byte) (TokenMaker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid symmetric key size: must be %d bytes, got %d bytes", chacha20poly1305.KeySize, len(symmetricKey))
	}
	pasetoMaker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: symmetricKey,
	}

	return pasetoMaker, nil
}

func (pm *PasetoMaker) CreateToken(userId uuid.UUID, identifier string, duration time.Duration) (string, error) {
	payload, err := NewPayload(userId, identifier, duration)
	if err != nil {
		return "", fmt.Errorf("failed to create payload: %w", err)
	}

	token, err := pm.paseto.Encrypt(pm.symmetricKey, payload, nil)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt token: %w", err)
	}
	return token, nil
}

func (pm *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	p := &Payload{}

	err := pm.paseto.Decrypt(token, pm.symmetricKey, p, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = p.Valid()
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return p, nil
}
