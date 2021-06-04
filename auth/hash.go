package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

const HashTokenBytes = 64

func GenerateHashToken() ([]byte, string, error) {
	b := make([]byte, HashTokenBytes)
	_, err := rand.Read(b)
	if err != nil {
		return nil, "", err
	}
	w := sha256.Sum256(b)
	return w[:], base64.URLEncoding.EncodeToString(b), nil
}

func RestoreHashToken(token string) ([]byte, error) {
	r, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}
	m := sha256.Sum256(r)
	return m[:], nil
}
