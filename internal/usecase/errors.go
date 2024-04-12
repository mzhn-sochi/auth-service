package usecase

import "errors"

var (
	ErrInvalidRole         = errors.New("invalid role")
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidToken        = errors.New("invalid token")
	ErrCannotGenerateToken = errors.New("cannot generate token")
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrPhoneAlreadyUsing   = errors.New("phone already using")
)
