package server

import (
	"context"
	"log/slog"
	"os"
	"smoothstart/auth"
	"smoothstart/handlers"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *SmoothStartServer) Routes() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	s.Server.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
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
	userHandler := handlers.NewUserHandler(s.DB)
	adminHandler := handlers.NewAdminHandler(s.DB)
	loginHandler := handlers.NewLoginHandler(s.DB, s.Redis, os.Getenv("SECRET_KEY"))
	auth := auth.NewAuthHandler(s.DB, s.Redis, os.Getenv("SECRET_KEY"))

	s.Server.Static("/", "assets")

	// index login page
	s.Server.GET("/", loginHandler.HandleLoginPage)
	s.Server.POST("/login", loginHandler.HandleLogin)
	s.Server.GET("/validate", auth.Validate(loginHandler.Validate))

	// user routes
	user := s.Server.Group("/user")
	user.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(os.Getenv("SECRET_KEY")),
		TokenLookup: "cookie:jwt",
	}))
	user.GET("/home", auth.Validate(userHandler.HomePage), auth.IsUser)
	user.GET("/plan", auth.Validate(userHandler.PlanPage), auth.IsUser)
	user.PUT("/plan", auth.Validate(userHandler.HandleUpdateStepStatus), auth.IsUser)
	user.POST("/logout", auth.Validate(loginHandler.HandleLogout), auth.IsUser)

	// admin routes
	admin := s.Server.Group("/admin")
	admin.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(os.Getenv("SECRET_KEY")),
		TokenLookup: "cookie:jwt",
	}))
	admin.GET("/home", auth.Validate(adminHandler.HomePage), auth.IsAdmin)
	admin.GET("/team", auth.Validate(adminHandler.TeamPage), auth.IsAdmin)
	admin.GET("/plans", auth.Validate(adminHandler.PlansPage), auth.IsAdmin)
	admin.GET("/plans/:id", auth.Validate(adminHandler.PlanPage), auth.IsAdmin)
	admin.GET("/plans/member/:id", auth.Validate(adminHandler.MemberPlanPage), auth.IsAdmin)
	admin.PUT("/plans/member/:id", auth.Validate(adminHandler.MemberPlanStepStatusUpdate), auth.IsAdmin)
	admin.POST("/logout", auth.Validate(loginHandler.HandleLogout), auth.IsAdmin)

}
