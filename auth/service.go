package auth

import (
	"github.com/thecactusblue/hana-id/model"
	"gorm.io/gorm"
	"strings"
)

type Service struct {
	db         *gorm.DB
	jwtFactory *JWTFactory
}

func NewService(db *gorm.DB, jwtFactory *JWTFactory) *Service {
	return &Service{db: db, jwtFactory: jwtFactory}
}

func CreateRefreshToken(user *model.User) (*model.RefreshToken, string, error) {
	t := &model.RefreshToken{
		UserID: user.ID,
	}
	s, err := t.Regen()
	if err != nil {
		return nil, "", err
	}
	return t, s, nil
}

func (s *Service) createRefresh(user *model.User) (string, error) {
	tok := &model.RefreshToken{
		UserID: user.ID,
	}
	rf, err := tok.Regen()
	if err != nil {
		return "", err
	}
	s.db.Create(tok)
	return rf, err
}

func (s *Service) Register(email, name, password string) error {
	u := &model.User{
		Email: strings.ToLower(strings.TrimSpace(email)),
		Name:  name,
	}
	u.SetPassword(password)
	s.db.Create(u)
	return nil
}

func (s *Service) Refresh(refreshToken string) (*AccessRefreshPair, error) {
	decoded, err := model.DecodeRefresh(refreshToken)
	if err != nil {
		return nil, err
	}
	rf := new(model.RefreshToken)
	err = s.db.Preload("User").First(rf, "token = ?", decoded).Error
	if err != nil {
		return nil, err
	}
	jwt, err := s.jwtFactory.IssueJWT(rf.User)
	if err != nil {
		return nil, err
	}
	r, err := rf.Regen()
	if err != nil {
		return nil, err
	}
	if err = s.db.Save(rf).Error; err != nil {
		return nil, err
	}
	return &AccessRefreshPair{
		AccessToken:  jwt,
		RefreshToken: r,
	}, nil
}

func (s *Service) Login(email, password string) (*AccessRefreshPair, error) {
	user := new(model.User)
	err := s.db.First(user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	if err := user.CheckPassword(password); err != nil {
		return nil, err
	}
	refresh, err := s.createRefresh(user)
	if err != nil {
		return nil, err
	}
	jwt, err := s.jwtFactory.IssueJWT(user)
	if err != nil {
		return nil, err
	}
	return &AccessRefreshPair{
		AccessToken:  jwt,
		RefreshToken: refresh,
	}, nil
}
