package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/mzhn-sochi/auth-service/internal/entity"
)

type claims struct {
	entity.UserClaims
	jwt.RegisteredClaims
}
