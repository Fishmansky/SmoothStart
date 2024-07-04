package auth

import (
	"database/sql"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	db    *sql.DB
	redis *redis.Client
}

func NewAuthHandler(s *sql.DB, r *redis.Client) *AuthHandler {
	return &AuthHandler{
		db:    s,
		redis: r,
	}
}
func (a AuthHandler) Validate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Validate if user is logged in
		_, err := c.Cookie("username")
		if err != nil {
			return err
		}
		return next(c)
	}
}

func (a AuthHandler) IsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Validate if user is logged in
		coockie, err := c.Cookie("username")
		if err != nil {
			return err
		}
		var isAdmin bool
		if err := a.db.QueryRow("SELECT is_admin FROM users WHERE username = $1", coockie.Value).Scan(&isAdmin); err != nil {
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
		coockie, err := c.Cookie("username")
		if err != nil {
			return err
		}
		var isAdmin bool
		if err := a.db.QueryRow("SELECT is_admin FROM users WHERE username = $1", coockie.Value).Scan(&isAdmin); err != nil {
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
