package auth

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"smoothstart/handlers"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func (a AuthHandler) ParseAccessToken(tokenString string) *handlers.JwtSSSToken {
	token, _ := jwt.ParseWithClaims(tokenString, &handlers.JwtSSSToken{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.secret), nil
	})
	return token.Claims.(*handlers.JwtSSSToken)
}

type AuthHandler struct {
	db     *sql.DB
	redis  *redis.Client
	secret string
}

func NewAuthHandler(sql *sql.DB, r *redis.Client, s string) *AuthHandler {
	return &AuthHandler{
		db:     sql,
		redis:  r,
		secret: s,
	}
}
func (a AuthHandler) verifyToken(tokenString string) error {
	token, err := jwt.ParseWithClaims(tokenString, &handlers.JwtSSSToken{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.secret), nil
	})
	if err != nil {
		return err
	}
	ctx := context.Background()
	result := a.redis.Get(ctx, tokenString)
	if result.Val() == "blacklisted" {
		return fmt.Errorf("Token blacklisted!")
	}
	ttl, err := a.redis.TTL(ctx, tokenString).Result()
	if err != nil {
		return err
	}
	if ttl == -2 {
		return fmt.Errorf("Token doesn't exist!")
	}
	if ttl == -1 {
		return fmt.Errorf("Token has no expiry date!")
	}
	if !token.Valid {
		return fmt.Errorf("Invalid token!")
	}

	return nil
}

func (a AuthHandler) Validate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Validate if user is logged in
		cookie, err := c.Cookie("jwt")
		if err != nil {
			return err
		}
		err = a.verifyToken(cookie.Value)
		if err != nil {
			return err
		}
		return next(c)
	}
}

func (a AuthHandler) IsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("jwt")
		if err != nil {
			return err
		}
		err = a.verifyToken(cookie.Value)
		if err != nil {
			return err
		}
		token := a.ParseAccessToken(cookie.Value)
		var isAdmin bool
		if err := a.db.QueryRow("SELECT is_admin FROM users WHERE username = $1", token.Username).Scan(&isAdmin); err != nil {
			if err == sql.ErrNoRows {
				return c.JSON(http.StatusUnauthorized, "User not found")
			}
			return c.JSON(http.StatusUnauthorized, err)
		}
		if !isAdmin {
			return c.JSON(http.StatusUnauthorized, "User is not admin")
		}

		return next(c)
	}
}

func (a AuthHandler) IsUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Validate if user is logged in
		cookie, err := c.Cookie("jwt")
		if err != nil {
			return err
		}
		err = a.verifyToken(cookie.Value)
		if err != nil {
			return err
		}
		token := a.ParseAccessToken(cookie.Value)
		var isAdmin bool
		if err := a.db.QueryRow("SELECT is_admin FROM users WHERE username = $1", token.Username).Scan(&isAdmin); err != nil {
			if err == sql.ErrNoRows {
				return c.JSON(http.StatusUnauthorized, "User not found")
			}
			return c.JSON(http.StatusUnauthorized, err)
		}
		if isAdmin {
			return c.JSON(http.StatusUnauthorized, "You are admin, not user")
		}

		return next(c)
	}
}
