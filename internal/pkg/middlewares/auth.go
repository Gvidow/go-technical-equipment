package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

var TokenJWTSecret = []byte("equipment")

type userID string

var ContextUserID userID = "user-id"

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tok := c.Request.Header.Get("Authorization")

		tokenJWT, err := jwt.Parse(tok, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("bad jwt token signing method: %s", token.Method.Alg())
			}
			return TokenJWTSecret, nil
		})
		if err != nil {
			return
		}

		if claims, ok := tokenJWT.Claims.(jwt.MapClaims); ok && tokenJWT.Valid {
			c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), ContextUserID, int(claims["id"].(float64))))
		}
	}
}

func RequireAuth() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		user := c.Request.Context().Value(ContextUserID)
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "требуется авторизация"})
			c.Abort()
		}
	})
}
