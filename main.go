package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"log/slog"
	"os"
	"smoothstart/auth"
	"smoothstart/handlers"
	"smoothstart/models"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	dbConf := os.Getenv("DBCONNSTR")
	DB, err = sql.Open("postgres", dbConf)
	if err != nil {
		log.Fatal(err)
	}
}

func DBInit() {
	_, err := DB.Query("DROP TABLE plans;")
	if err != nil {
		slog.Error(err.Error())
	}
	_, err = DB.Query("DROP TABLE users;")
	if err != nil {
		slog.Error(err.Error())
	}
	_, err = DB.Query("CREATE TABLE users (id SERIAL PRIMARY KEY, username VARCHAR(255) NOT NULL UNIQUE, fname VARCHAR(255) NOT NULL, sname VARCHAR(255) NOT NULL, password VARCHAR(255) NOT NULL, is_admin BOOLEAN NOT NULL);")
	if err != nil {
		slog.Error(err.Error())
	}
	_, err = DB.Query("CREATE TABLE plans (id SERIAL PRIMARY KEY, name VARCHAR(255) NOT NULL, description TEXT, userid INT NOT NULL, steps JSON NULL, FOREIGN KEY (userid) REFERENCES users(id));")
	if err != nil {
		slog.Error(err.Error())
	}
	_, err = DB.Query("INSERT INTO users (username, fname, sname, password, is_admin) VALUES ($1, $2, $3, $4, $5);", "admin", "ad", "min", "test123", 1)
	if err != nil {
		slog.Error(err.Error())
	}
	_, err = DB.Query("INSERT INTO users (username, fname, sname, password, is_admin) VALUES ($1, $2, $3, $4, $5);", "user1", "John", "Long", "test123", 0)
	if err != nil {
		slog.Error(err.Error())
	}
	steps := []models.Step{{ID: 0, Description: "Create account", Done: false}, {ID: 1, Description: "Create password", Done: false}}
	d, _ := json.Marshal(steps)
	_, err = DB.Query("INSERT INTO plans (name, description, userid, steps) VALUES ($1, $2, $3, $4);", "user1 plan", "plan for user1", 2, d)
	if err != nil {
		slog.Error(err.Error())
	}
}

func main() {
	ConnectDB()
	DBInit()
	e := echo.New()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))
	// handlers
	uH := handlers.NewUserHandler(DB)
	aH := handlers.NewAdminHandler(DB)
	lH := handlers.NewLoginHandler(DB)
	auth := auth.AuthHandler{DB: DB}

	e.Static("/", "assets")

	// index login page
	e.GET("/", lH.HandleLoginPage)
	e.POST("/login", lH.HandleLogin)

	// user routes
	user := e.Group("/user")
	user.GET("/home", auth.Validate(uH.HomePage), auth.IsUser)
	user.GET("/plans", auth.Validate(uH.PlansPage), auth.IsUser)
	user.GET("/plans/:id", auth.Validate(uH.HandlePlan), auth.IsUser)

	// admin routes
	admin := e.Group("/admin")
	admin.GET("/home", auth.Validate(aH.HomePage), auth.IsAdmin)

	e.Logger.Fatal(e.Start(":8080"))
}
