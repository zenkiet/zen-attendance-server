package token

import (
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
)

type Maker struct {
	symmetricKey paseto.V4SymmetricKey
	implicit     []byte
}

var (
	ErrInvalidKeyLength = fmt.Errorf("invalid key length: must be 64 hex characters (32 bytes)")
	ErrInvalidImplicit  = fmt.Errorf("invalid implicit value: must be non-empty")
	ErrInvalidSubject   = fmt.Errorf("invalid subject: must be non-empty")
	ErrInvalidDuration  = fmt.Errorf("invalid duration: must be greater than zero")
	ErrInvalidToken     = fmt.Errorf("invalid token")
	ErrTokenExpired     = fmt.Errorf("token has expired")
	ErrInvalidDeviceID  = fmt.Errorf("invalid device ID: must be non-empty")
	ErrInvalidNonce     = fmt.Errorf("invalid nonce: must be non-empty")
)

func NewMaker(hexKey string, implicit []byte) (*Maker, error) {
	if len(implicit) == 0 {
		return nil, ErrInvalidKeyLength
	}

	if len(hexKey) != 64 { // 64 === 32 bytes in hex
		return nil, ErrInvalidKeyLength
	}

	if len(implicit) == 0 {
		return nil, ErrInvalidImplicit
	}

	key, err := paseto.V4SymmetricKeyFromHex(hexKey)
	if err != nil {
		return nil, fmt.Errorf("could not create symmetric key from hex string: %w", err)
	}

	return &Maker{
		symmetricKey: key,
		implicit:     implicit,
	}, nil
}

func (maker *Maker) CreateToken(subject string, duration time.Duration) (string, error) {
	if len(subject) == 0 {
		return "", ErrInvalidSubject
	}

	if duration <= 0 {
		return "", ErrInvalidDuration
	}

	token := paseto.NewToken()
	now := time.Now().UTC()

	token.SetJti(uuid.NewString())
	token.SetSubject(subject)
	token.SetIssuedAt(now)
	token.SetNotBefore(now)
	token.SetIssuedAt(now)
	token.SetExpiration(now.Add(duration))

	encrypted := token.V4Encrypt(maker.symmetricKey, maker.implicit)

	return encrypted, nil
}

func (maker *Maker) VerifyToken(tokenString string) (*paseto.Token, error) {
	if len(tokenString) == 0 {
		return nil, ErrInvalidToken
	}
	parser := paseto.NewParser()
	now := time.Now().UTC()

	parser.AddRule(paseto.NotExpired())
	parser.AddRule(paseto.ValidAt(now))

	tok, err := parser.ParseV4Local(maker.symmetricKey, tokenString, maker.implicit)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return tok, nil
}

func (maker *Maker) GetTokenExpiration(tokenString string) (time.Time, error) {
	token, err := maker.VerifyToken(tokenString)
	if err != nil {
		return time.Time{}, err
	}

	expiration, err := token.GetExpiration()
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get expiration: %w", err)
	}

	return expiration, nil
}
