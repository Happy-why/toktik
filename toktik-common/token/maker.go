package token

import (
	"errors"
	"time"
)

const minSecretKeySize = 32

var ErrSecretLen = errors.New("密钥长度不合法")

// Maker is an interface for managing tokens
type Maker interface {
	// CreateToken creates a new token for a specific username and duration
	CreateToken(content []byte, duration time.Duration) (string, *Payload, error)
	// VerifyToken checks if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}
