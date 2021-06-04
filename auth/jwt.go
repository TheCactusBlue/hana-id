package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/thecactusblue/hana-id/model"
	"time"
)

type JWTFactory struct {
	secret []byte
}

func NewJWTFactory(secret string) *JWTFactory {
	return &JWTFactory{secret: []byte(secret)}
}

type AccessRefreshPair struct {
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}

type Claims struct {
	jwt.StandardClaims
	Name string `json:"name"`
}

func (f *JWTFactory) IssueJWT(u *model.User) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		Name: u.Name,
		StandardClaims: jwt.StandardClaims{
			Subject:   u.ID.String(),
			Issuer:    "multitudehq.com",
			ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		},
	}).SignedString(f.secret)
}
