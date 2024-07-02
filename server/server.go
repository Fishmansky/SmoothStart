package server

import (
	"database/sql"
	"encoding/json"
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

func (s *SmoothStartServer) ConnectDB(psqlDSN string) {
	var err error
	s.DB, err = sql.Open("postgres", psqlDSN)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *SmoothStartServer) ConnectRedis(redisConf string) {
	ops, err := redis.ParseURL(redisConf)
	if err != nil {
		log.Fatal(err)
	}
	s.Redis = redis.NewClient(ops)
}

func (s *SmoothStartServer) DBInit() {
	_, err := s.DB.Query("DROP TABLE plan_templates;")
	if err != nil {
		slog.Error(err.Error())
	}
	_, err = s.DB.Query("DROP TABLE plans;")
	if err != nil {
		slog.Error(err.Error())
	}
	_, err = s.DB.Query("DROP TABLE users;")
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
	psqlDSN := os.Getenv("PSQL_DSN")
	redisConf := os.Getenv("REDIS_CONF")
	s.ConnectDB(psqlDSN)
	s.ConnectRedis(redisConf)
	s.DBInit()
	s.Routes()
	return s
}
