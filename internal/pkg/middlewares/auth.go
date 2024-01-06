package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"

	"github.com/gvidow/go-technical-equipment/internal/app/config"
	"github.com/gvidow/go-technical-equipment/internal/app/ds"
	"github.com/gvidow/go-technical-equipment/internal/app/redis"
)

type user string

var ContextUser user = "user-with-role"

type UserWithRole struct {
	UserID int
	Role   ds.Role
}

func Auth(cfg config.JWTConfig, client *redis.Client) gin.HandlerFunc {
	var jwtPrefix = cfg.TokenType + " "

	return func(c *gin.Context) {
		tok := c.Request.Header.Get(cfg.Header)

		if !strings.HasPrefix(tok, jwtPrefix) {
			return
		}

		jwtStr := tok[len(jwtPrefix):]
		ok, err := client.CheckJWTInBlackList(c, jwtStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "не удалось аутентифицировать пользователя"})
			c.Abort()
			return
		}
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": "недействительный токен"})
			c.Abort()
			return
		}

		claims := &ds.JWTClaims{}

		tokenJWT, err := jwt.ParseWithClaims(jwtStr, claims, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("bad jwt token signing method: %s", token.Method.Alg())
			}
			return cfg.SecretToken, nil
		})
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": "недействительный токен"})
			c.Abort()
			return
		}

		if tokenJWT.Valid {
			c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), ContextUser, UserWithRole{
				UserID: claims.UserID,
				Role:   claims.Role,
			}))
		}
	}
}

func RequireAuth(requiredRoles ...ds.Role) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		user := c.Request.Context().Value(ContextUser)
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "требуется авторизация"})
			c.Abort()
			return
		}
		if len(requiredRoles) == 0 {
			return
		}
		userRole := user.(UserWithRole).Role
		for _, role := range requiredRoles {
			if userRole == role {
				return
			}
		}
		c.JSON(http.StatusMethodNotAllowed, gin.H{"status": "error", "message": "пользователь не имеет прав на совершение данного действия"})
		c.Abort()
	})
}
