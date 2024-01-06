package ds

import (
	jwt "github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	jwt.RegisteredClaims

	UserID int
	Role   Role
}
