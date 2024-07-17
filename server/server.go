package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"smoothstart/models"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

type SmoothStartServer struct {
	Server *echo.Echo
	DB     *sql.DB
	Redis  *redis.Client
}

func LoadEnvs() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func (s *SmoothStartServer) ConnectDB() {
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	db := os.Getenv("POSTGRES_DB")
	pass := os.Getenv("POSTGRES_PASSWORD")
	var err error
	connStr := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", user, pass, host, db)
	s.DB, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
}

func (s *SmoothStartServer) ConnectRedis() {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	ops := &redis.Options{
		Addr:     host + ":" + port,
		Password: "",
		DB:       0,
	}
	s.Redis = redis.NewClient(ops)
	ctx := context.Background()
	if err := s.Redis.Ping(ctx).Err(); err != nil {
		log.Println(err)
	}
}

func (s *SmoothStartServer) DBInit() {
	_, err := s.DB.Query("DROP TABLE IF EXISTS plan_templates;")
	if err != nil {
		slog.Error(err.Error())
	}
	_, err = s.DB.Query("DROP TABLE IF EXISTS plans;")
	if err != nil {
		slog.Error(err.Error())
	}
	_, err = s.DB.Query("DROP TABLE IF EXISTS users;")
	if err != nil {
		slog.Error(err.Error())
	}
	_, err = s.DB.Query("CREATE TABLE users (id SERIAL PRIMARY KEY, username VARCHAR(255) NOT NULL UNIQUE, fname VARCHAR(255) NOT NULL, sname VARCHAR(255) NOT NULL, password VARCHAR(255) NOT NULL, is_admin BOOLEAN NOT NULL);")
	if err != nil {
		slog.Error(err.Error())
	}
	_, err = s.DB.Query("CREATE TABLE plans (id SERIAL PRIMARY KEY, name VARCHAR(255) NOT NULL, description TEXT, steps JSON NULL, user_id INT NOT NULL, FOREIGN KEY (user_id) REFERENCES users(id));")
	if err != nil {
		slog.Error(err.Error())
	}
	_, err = s.DB.Query("CREATE TABLE plan_templates (id SERIAL PRIMARY KEY, name VARCHAR(255) NOT NULL, description TEXT, steps JSON NULL);")
	if err != nil {
		slog.Error(err.Error())
	}
	_, err = s.DB.Query("INSERT INTO users (username, fname, sname, password, is_admin) VALUES ($1, $2, $3, $4, $5);", "admin", "ad", "min", "test123", 1)
	if err != nil {
		slog.Error(err.Error())
	}
	_, err = s.DB.Query("INSERT INTO users (username, fname, sname, password, is_admin) VALUES ($1, $2, $3, $4, $5);", "user1", "John", "Long", "test123", 0)
	if err != nil {
		slog.Error(err.Error())
	}
	steps := []models.Step{{ID: 0, Description: "Create account", Done: false}, {ID: 1, Description: "Create password", Done: false}}
	d, _ := json.Marshal(steps)
	_, err = s.DB.Query("INSERT INTO plans (name, description, steps, user_id) VALUES ($1, $2, $3, $4);", "John plan", "plan for John", d, 2)
	if err != nil {
		slog.Error(err.Error())
	}
	steps2 := []models.Step{{ID: 0, Description: "Create account", Done: false}, {ID: 1, Description: "Create password", Done: false}}
	d2, _ := json.Marshal(steps2)
	_, err = s.DB.Query("INSERT INTO plan_templates (name, description, steps) VALUES ($1, $2, $3);", "Test plan", "plan", d2)
	if err != nil {
		slog.Error(err.Error())
	}
}
func NewSSS() *SmoothStartServer {
	s := &SmoothStartServer{
		Server: echo.New(),
	}
	LoadEnvs()
	s.ConnectDB()
	s.ConnectRedis()
	s.DBInit()
	s.Routes()
	return s
}
