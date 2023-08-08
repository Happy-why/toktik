package token

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalided")
)

type User struct { // 伪造一个Payload的Content
	UserID   int64  `json:"userid"`
	UserName string `json:"username"`
}

// Payload contains the payload data of the token
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Content   []byte    `json:"content"` // 可以是任意内容
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(content []byte, duration time.Duration) (*Payload, error) {
	// NewRandom return a random uuid.UUID
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		Content:   content,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
