package model

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Model struct {
	ID uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
}

type User struct {
	Model
	Name     string  `gorm:"unique"`
	Email    string  `gorm:"unique"`
	Password []byte  `json:"-"`
	Roles    []*Role `json:"-"`
}

func (u *User) SetPassword(password string) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		panic(err)
	}
	u.Password = []byte(hash)
}

var ErrIncorrectPassword error = echo.NewHTTPError(400, echo.Map{
	"code": "IncorrectPassword",
})

func (u *User) CheckPassword(password string) error {
	succ, err := argon2id.ComparePasswordAndHash(password, string(u.Password))
	if err != nil {
		return err
	} else if !succ {
		return ErrIncorrectPassword
	}
	return nil
}

type Role struct {
	Model
	UserID uuid.UUID `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE;"`
	Name   string    `gorm:"not null"`
}

type RefreshToken struct {
	Model

	UserID uuid.UUID `gorm:"type:uuid;not null;index;constraint:OnDelete:CASCADE;"`
	User   *User
	Token  []byte `gorm:"unique;not null"`
}

func (r *RefreshToken) Regen() (string, error) {
	b := make([]byte, 64)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	w := sha256.Sum256(b)
	r.Token = w[:]
	return base64.URLEncoding.EncodeToString(b), nil
}

func DecodeRefresh(token string) ([]byte, error) {
	r, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}
	m := sha256.Sum256(r)
	return m[:], nil
}

// identities
type Credential interface{}
