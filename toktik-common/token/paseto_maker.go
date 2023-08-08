package token

import (
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
	"time"
)

type PasetoMaker struct {
	paseto *paseto.V2
	Key    []byte //symmetric对称
}

func NewPasetoMaker(Key string) (Maker, error) {
	if len(Key) != chacha20poly1305.KeySize {
		return nil, ErrSecretLen
	}

	maker := &PasetoMaker{
		paseto: paseto.NewV2(),
		Key:    []byte(Key),
	}
	return maker, nil
}

// CreateToken 构造token
func (maker *PasetoMaker) CreateToken(content []byte, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(content, duration)
	if err != nil {
		return "", nil, err
	}

	token, err := maker.paseto.Encrypt(maker.Key, payload, nil)
	return token, payload, nil
}

// VerifyToken checks if the token is valid or not
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	// 验证token是否有效
	err := maker.paseto.Decrypt(token, maker.Key, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// 验证token是否过期
	err = payload.Valid()
	if err != nil {
		return nil, ErrExpiredToken
	}

	return payload, nil
}
