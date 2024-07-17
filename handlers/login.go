package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"html"
	"log"
	"net/http"
	"smoothstart/models"
	"smoothstart/views/layout"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type JwtSSSToken struct {
	ID       uint32 `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (l LoginHandler) createAccessToken(u *models.User) (string, error) {
	id := uuid.New().ID()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&JwtSSSToken{
			id,
			u.Username,
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
			},
		})

	tokenString, err := token.SignedString([]byte(l.secret))
	if err != nil {
		log.Println(err)
		return "", err
	}
	ctx := context.Background()
	err = l.redis.Set(ctx, tokenString, "valid", time.Minute*15).Err()
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (l LoginHandler) createRefreshToken(u *models.User) (string, error) {
	id := uuid.New().ID()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&JwtSSSToken{
			id,
			u.Username,
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			},
		})

	tokenString, err := token.SignedString([]byte(l.secret))
	if err != nil {
		log.Println(err)
		return "", err
	}

	return tokenString, nil
}

func (l LoginHandler) verifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return l.secret, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("Invalid token!")
	}

	return nil
}

func (l LoginHandler) invalidateToken(token string) error {
	ctx := context.Background()
	err := l.redis.Set(ctx, token, "blacklisted", 0).Err()
	if err != nil {
		return err
	}
	return nil
}

type LoginHandler struct {
	db     *sql.DB
	redis  *redis.Client
	secret string
}

func NewLoginHandler(d *sql.DB, r *redis.Client, s string) *LoginHandler {
	return &LoginHandler{
		db:     d,
		redis:  r,
		secret: s,
	}
}

func (l LoginHandler) getUser(username string, password string) (*models.User, error) {
	var user models.User
	if err := l.db.QueryRow("SELECT id, username, fname, sname, is_admin FROM users WHERE username = $1 AND password = $2", username, password).Scan(&user.ID, &user.Username, &user.Fname, &user.Sname, &user.IsAdmin); err != nil {
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
	//cookie, err := c.Cookie("jwt")
	//if err != nil {
	//	return err
	//}
	//var user int
	//if err := l.db.QueryRow("SELECT id FROM users WHERE username = $1", token.Username).Scan(&user); err != nil {
	//	if err == sql.ErrNoRows {
	//		return fmt.Errorf("%s\n", "User not found")
	//	}
	//	return fmt.Errorf("%s\n", err)
	//}
	return nil
}
func (l LoginHandler) HandleLogout(c echo.Context) error {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		return err
	}
	l.invalidateToken(cookie.Value)
	c.Response().Header().Set("HX-Location", "/")
	return c.JSON(http.StatusOK, "Logout succesfull")
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
	// set jwt
	tokenString, err := l.createAccessToken(u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "No username found")
	}
	cookie := new(http.Cookie)
	cookie.Name = "jwt"
	cookie.Value = tokenString
	cookie.HttpOnly = true
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)
	if u.IsAdmin {
		c.Response().Header().Set("HX-Location", "/admin/home")
	} else {
		c.Response().Header().Set("HX-Location", "/user/home")
	}
	return c.JSON(http.StatusOK, "Login succesfull")
}
