package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"html"
	"net/http"
	"smoothstart/models"
	"smoothstart/views/layout"
	"time"

	"github.com/labstack/echo/v4"
)

type LoginHandler struct {
	db *sql.DB
}

func NewLoginHandler(d *sql.DB) *LoginHandler {
	return &LoginHandler{
		db: d,
	}
}

func (l LoginHandler) getUser(username string, password string) (*models.User, error) {
	var user models.User
	if err := l.db.QueryRow("SELECT id, username, fname, sname is_admin FROM users WHERE username = $1 AND password = $2", username, password).Scan(&user); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s\n", "User not found")
		}
		return nil, fmt.Errorf("%s\n", err)
	}
	return &user, nil
}

func ValidateLoginData(data *models.LoginData) error {
	// validate login data
	if data.Username == "" || data.Password == "" {
		return errors.New("")
	}
	// validated
	return nil
}

func (l LoginHandler) HandleLoginPage(c echo.Context) error {
	return render(c, layout.LoginPage(models.LoginData{}))
}

func (l LoginHandler) Validate(c echo.Context) error {
	cookie, err := c.Cookie("username")
	if err != nil {
		return err
	}
	var user int
	if err := l.db.QueryRow("SELECT id FROM users WHERE username = $1", cookie.Value).Scan(&user); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("%s\n", "User not found")
		}
		return fmt.Errorf("%s\n", err)
	}
	return nil
}
func (l LoginHandler) HandleLogin(c echo.Context) error {
	var input models.LoginData
	if err := c.Bind(&input); err != nil {
		return c.String(http.StatusBadRequest, "Bad request")
	}
	data := models.LoginData{}
	data.Username = html.EscapeString(input.Username)
	data.Password = html.EscapeString(input.Password)
	if err := ValidateLoginData(&data); err != nil {
		return c.String(http.StatusBadRequest, "Bad login data")
	}
	// find user in DB
	u, err := l.getUser(data.Username, data.Password)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// set coockie
	cookie := new(http.Cookie)
	cookie.Name = "username"
	cookie.Value = u.Username
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)
	if u.IsAdmin {
		c.Response().Header().Set("HX-Location", "/admin/home")
	} else {
		c.Response().Header().Set("HX-Location", "/user/home")
	}
	return c.JSON(http.StatusOK, "Login succesfull")
}
