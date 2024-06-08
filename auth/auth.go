package auth

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	DB *sql.DB
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
		if err := a.DB.QueryRow("SELECT is_admin FROM users WHERE username = $1", coockie.Value).Scan(&isAdmin); err != nil {
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
