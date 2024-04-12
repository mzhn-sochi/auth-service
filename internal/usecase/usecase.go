package usecase

import (
	"fmt"
	"github.com/mzhn-sochi/auth-service/internal/config"
	"github.com/mzhn-sochi/auth-service/internal/entity"
	"github.com/mzhn-sochi/auth-service/internal/lib/jwt"
	"log"
	"time"
)

type UserStorage interface {
	Get(id string) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	GetByPhone(phone string) (*entity.User, error)

	Create(user *entity.User) error
	Update(user *entity.User) error
	Delete(id int) error
}

type TokenStorage interface {
	Get(userId string) (string, error)
	Save(userId string, token string) error
	Delete(userId string) error
}

type UseCase struct {
	userStorage  UserStorage
	tokenStorage TokenStorage
	cfg          *config.Config
}

func New(cfg *config.Config, userStorage UserStorage, tokenStorage TokenStorage) *UseCase {
	return &UseCase{
		userStorage:  userStorage,
		tokenStorage: tokenStorage,
		cfg:          cfg,
	}
}

func (s *UseCase) SignUp(user *entity.User) (*entity.Tokens, error) {

	if _, err := s.userStorage.GetByPhone(user.Phone); err == nil {
		// TODO handle errors
		return nil, fmt.Errorf("phone %s is already in use", user.Phone)
	}

	tokens, err := s.generateJwtPair(user.GetClaims())
	if err != nil {
		return nil, fmt.Errorf("failed to generate jwt pair: %v", err)
	}

	if err := s.tokenStorage.Save(user.Id, tokens.Refresh); err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %v", err)
	}

	if err := s.userStorage.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	return tokens, nil
}

func (s *UseCase) SignIn(user *entity.User) (*entity.Tokens, error) {
	u, err := s.userStorage.GetByPhone(user.Phone)
	if err != nil {
		// TODO handle errors
		return nil, fmt.Errorf("user with phone %s not found", user.Email)
	}

	tokens, err := s.generateJwtPair(u.GetClaims())
	if err != nil {
		return nil, fmt.Errorf("failed to generate jwt pair: %v", err)
	}

	if err := s.tokenStorage.Save(u.Id, tokens.Refresh); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %v", err)
	}

	return tokens, nil
}

func (s *UseCase) SignOut(accessToken string) error {
	claims, err := jwt.Validate(accessToken, s.cfg.JWT.Access.Secret)
	if err != nil {
		return fmt.Errorf("failed to validate access token: %v", err)
	}

	return s.tokenStorage.Delete(claims.Id)
}

func (s *UseCase) Authenticate(accessToken string, role entity.Role) error {
	claims, err := jwt.Validate(accessToken, s.cfg.JWT.Access.Secret)
	if err != nil {
		log.Printf("failed to validate access token: %v", err)
		return ErrInvalidToken
	}

	u, err := s.userStorage.Get(claims.Id)
	if err != nil {
		log.Printf("failed to get user: %v", err)
		return ErrUserNotFound
	}

	if role != u.Role {
		log.Printf("user %s has role %s, but %s is required", u.Id, u.Role, role)
		return ErrInvalidRole
	}

	return nil
}

func (s *UseCase) Refresh(refreshToken string) (*entity.Tokens, error) {

	claims, err := jwt.Validate(refreshToken, s.cfg.JWT.Refresh.Secret)
	if err != nil {
		log.Printf("failed to validate refresh token: %v", err)
		return nil, ErrInvalidToken
	}

	tokens, err := s.generateJwtPair(claims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate jwt pair: %v", err)
	}

	if err := s.tokenStorage.Save(claims.Id, tokens.Refresh); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %v", err)
	}

	return tokens, nil
}

func (s *UseCase) generateJwtPair(claims *entity.UserClaims) (*entity.Tokens, error) {

	refresh, err := jwt.Generate(claims, time.Duration(s.cfg.JWT.Refresh.TTL)*time.Hour, []byte(s.cfg.JWT.Refresh.Secret))
	if err != nil {
		return nil, err
	}

	access, err := jwt.Generate(claims, time.Duration(s.cfg.JWT.Access.TTL)*time.Hour, []byte(s.cfg.JWT.Access.Secret))
	if err != nil {
		return nil, err
	}

	return &entity.Tokens{
		Access:  access,
		Refresh: refresh,
	}, nil
}
