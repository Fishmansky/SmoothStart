package handlers

import (
	"database/sql"
	"smoothstart/models"
	"smoothstart/views/user"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	db *sql.DB
}

func NewUserHandler(d *sql.DB) *UserHandler {
	return &UserHandler{
		db: d,
	}
}

func (u UserHandler) HomePage(c echo.Context) error {
	return render(c, user.Home())
}

func (u UserHandler) PlansPage(c echo.Context) error {
	P := []models.Plan{{ID: 0, Name: "Admin", Description: "Route for administrators", Steps: []models.Step{{ID: 0, Description: "Create password", Done: true}, {ID: 1, Description: "Create account"}}}, {ID: 1, Name: "Sales", Description: "Route for salesman", Steps: []models.Step{{ID: 0, Description: "Create password"}, {ID: 1, Description: "Create account"}}}}
	return render(c, user.PlansGrid(P))
}

func (u UserHandler) HandlePlan(c echo.Context) error {
	id := c.Param("id")
	planId, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	p, err := models.GetUserPlan(planId)
	if err != nil {
		return err
	}
	return render(c, user.Plan(p))
}
