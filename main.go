package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"smoothstart/models"
	"smoothstart/views"
	"strconv"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var P []models.Plan

func renderView(c echo.Context, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return cmp.Render(c.Request().Context(), c.Response().Writer)
}
func indexView(c echo.Context) error {
	cmp := views.IndexComponent()
	return renderView(c, cmp)
}

func homeView(c echo.Context) error {
	v := views.HomeComponent()
	return renderView(c, v)
}
func teamView(c echo.Context) error {
	v := views.TeamComponent()
	return renderView(c, v)
}
func plansView(c echo.Context) error {
	v := views.PlansGrid(P)
	return renderView(c, v)
}
func planView(c echo.Context) error {
	id := c.Param("id")
	pID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Plan ID conversion failed")
	}
	plan := views.Plan(P[pID])
	return renderView(c, plan)
}

func main() {
	P = []models.Plan{{ID: 0, Name: "Admin", Description: "Route for administrators", Steps: []models.Step{{ID: 0, Description: "Create password", Done: true}, {ID: 1, Description: "Create account"}}}, {ID: 1, Name: "Sales", Description: "Route for salesman", Steps: []models.Step{{ID: 0, Description: "Create password"}, {ID: 1, Description: "Create account"}}}}
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
	e.Static("/", "assets")
	e.GET("/", indexView)
	e.GET("/plans", plansView)
	e.GET("/plans/:id", planView)
	e.GET("/home", homeView)
	e.GET("/team", teamView)
	e.Logger.Fatal(e.Start(":8080"))
}
