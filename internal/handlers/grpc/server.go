package grpc

import (
	"context"
	"errors"
	"github.com/mzhn-sochi/auth-service/api/auth"
	"github.com/mzhn-sochi/auth-service/internal/config"
	"github.com/mzhn-sochi/auth-service/internal/entity"
	"github.com/mzhn-sochi/auth-service/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

var _ auth.AuthServer = (*Server)(nil)

type AuthUseCase interface {
	SignUp(ctx context.Context, user *entity.User) (*entity.Tokens, error)
	SignIn(ctx context.Context, user *entity.User) (*entity.Tokens, error)
	SingOut(ctx context.Context, accessToken string) error
	Authenticate(ctx context.Context, accessToken string, role entity.Role) error
	Refresh(ctx context.Context, refreshToken string) (*entity.Tokens, error)
}

type Server struct {
	cfg *config.Config

	uc AuthUseCase

	auth.UnimplementedAuthServer
}

func New(cfg *config.Config, uc AuthUseCase) *Server {
	return &Server{
		cfg: cfg,
		uc:  uc,
	}
}

func (s *Server) SignIn(ctx context.Context, request *auth.SignInRequest) (*auth.Tokens, error) {
	user := &entity.User{
		UserCredentials: entity.UserCredentials{
			Phone:    request.Phone,
			Password: request.Password,
		},
	}

	tokens, err := s.uc.SignIn(ctx, user)
	if err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		if errors.Is(err, usecase.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &auth.Tokens{
		Access:  tokens.Access,
		Refresh: tokens.Refresh,
	}, nil
}

func (s *Server) SignUp(ctx context.Context, request *auth.SignUpRequest) (*auth.Tokens, error) {
	user := &entity.User{
		UserCredentials: entity.UserCredentials{
			Phone:    request.Phone,
			Password: request.Password,
		},
	}

	tokens, err := s.uc.SignUp(ctx, user)
	if err != nil {
		if errors.Is(err, usecase.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &auth.Tokens{
		Access:  tokens.Access,
		Refresh: tokens.Refresh,
	}, nil
}

func (s *Server) SignOut(ctx context.Context, request *auth.SignOutRequest) (*auth.Empty, error) {

	err := s.uc.SingOut(ctx, request.AccessToken)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidToken) {
			return nil, status.Error(codes.InvalidArgument, "invalid token")
		}

		if errors.Is(err, usecase.ErrSessionNotFound) {
			return nil, status.Error(codes.NotFound, "session not found")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &auth.Empty{}, nil
}

func (s *Server) Auth(ctx context.Context, request *auth.AuthRequest) (*auth.Empty, error) {

	if err := s.uc.Authenticate(ctx, request.AccessToken, entity.Role(request.Role)); err != nil {
		if errors.Is(err, usecase.ErrTokenExpired) {
			return nil, status.Error(codes.Unauthenticated, "token expired")
		}

		if errors.Is(err, usecase.ErrInvalidRole) {
			return nil, status.Error(codes.PermissionDenied, "invalid role")
		}

		if errors.Is(err, usecase.ErrSessionNotFound) {
			return nil, status.Error(codes.NotFound, "session not found")
		}

		if errors.Is(err, usecase.ErrInvalidToken) {
			return nil, status.Error(codes.InvalidArgument, "invalid token")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &auth.Empty{}, nil
}

func (s *Server) Refresh(ctx context.Context, request *auth.RefreshRequest) (*auth.Tokens, error) {

	log := ctx.Value("logger").(*slog.Logger).With("method", "Refresh")

	tokens, err := s.uc.Refresh(ctx, request.RefreshToken)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidToken) {
			log.Debug("invalid token")
			return nil, status.Error(codes.InvalidArgument, "invalid token")
		}

		if errors.Is(err, usecase.ErrSessionNotFound) {
			log.Debug("session not found")
			return nil, status.Error(codes.NotFound, "session not found")
		}

		log.Debug("internal server error", slog.String("err", err.Error()))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &auth.Tokens{
		Access:  tokens.Access,
		Refresh: tokens.Refresh,
	}, nil
}
