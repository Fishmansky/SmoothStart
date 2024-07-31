package models

import (
	"database/sql"
	"fmt"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Fname    string `json:"fname"`
	Lname    string `json:"lname"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
	HasPlan  bool
}

func FindUser(db *sql.DB, data LoginData) (User, error) {
	var u User
	if err := db.QueryRow("SELECT id, username, is_admin FROM users WHERE username = $1 AND password = $2", data.Username, data.Password).Scan(&u.ID, &u.Username, &u.IsAdmin); err != nil {
		if err == sql.ErrNoRows {
			return u, fmt.Errorf("%s\n", "User not found")
		}
		return u, fmt.Errorf("%s\n", err)
	}
	return u, nil
}

func GetUserID(db *sql.DB, username string) (int, error) {
	var uID int
	if err := db.QueryRow("SELECT id  FROM users WHERE username = $1", username).Scan(&uID); err != nil {
		if err == sql.ErrNoRows {
			return uID, fmt.Errorf("%s\n", "User not found")
		}
		return uID, fmt.Errorf("%s\n", err)
	}
	return uID, nil
}
