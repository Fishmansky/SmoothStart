package handlers

import (
	"database/sql"
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

func ValidateLoginData(data *models.LoginData) error {
	// validate login data
	if data.Password == "" || data.Username == "" {
		return fmt.Errorf("%s\n", "Password or username is blank")
	}
	data.Username = html.EscapeString(data.Username)
	data.Password = html.EscapeString(data.Password)
	return nil
}

func (l LoginHandler) HandleLoginPage(c echo.Context) error {
	return render(c, layout.LoginPage(models.LoginData{}))
}
func (l LoginHandler) HandleLogin(c echo.Context) error {
	var data models.LoginData
	if err := c.Bind(&data); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	if err := ValidateLoginData(&data); err != nil {
		return err
	}
	// find user in DB
	u, err := models.FindUser(l.db, data)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// set coockie
	cookie := new(http.Cookie)
	cookie.Name = "username"
	cookie.Value = u.Username
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, "Login succesfull")
}
