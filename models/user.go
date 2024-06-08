package models

import (
	"database/sql"
	"fmt"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Fname    string `json:"fname"`
	Sname    string `json:"sname"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}

func GetUserPlan(id int) (Plan, error) {
	return Plan{}, nil
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
