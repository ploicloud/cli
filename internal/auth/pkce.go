package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

type PKCE struct {
	Verifier  string
	Challenge string
	State     string
}

func NewPKCE() (*PKCE, error) {
	v, err := randB64URL(32)
	if err != nil {
		return nil, err
	}
	s, err := randB64URL(16)
	if err != nil {
		return nil, err
	}
	sum := sha256.Sum256([]byte(v))
	return &PKCE{
		Verifier:  v,
		Challenge: base64.RawURLEncoding.EncodeToString(sum[:]),
		State:     s,
	}, nil
}

func randB64URL(n int) (string, error) {
	if n <= 0 {
		return "", errors.New("invalid length")
	}
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
