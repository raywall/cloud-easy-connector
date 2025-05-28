package auth

import (
	tm "github.com/raywall/cloud-easy-connector/internal/auth"
)

type AutoManagedToken interface {
	Start() error
	GetToken() (string, error)
	Stop()
	RefreshLoop()
	RefreshToken() error
}

type AuthRequest tm.AuthRequest

// NewAutoManagedToken cria uma nova instância de TokenManager
func NewAutoManagedToken(apiURL string, authRequest AuthRequest, certSkipVerify bool) AutoManagedToken {
	return tm.NewManagedToken(apiURL, tm.AuthRequest(authRequest), certSkipVerify)
}

