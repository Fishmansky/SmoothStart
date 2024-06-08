package handlers

import (
	"database/sql"
	"smoothstart/models"
	"smoothstart/views/admin"

	"github.com/labstack/echo/v4"
)

type AdminHandler struct {
	db *sql.DB
}

func NewAdminHandler(d *sql.DB) *AdminHandler {
	return &AdminHandler{
		db: d,
	}
}

func (a AdminHandler) HomePage(ctx echo.Context) error {
	return render(ctx, admin.Home())
}
func (a AdminHandler) TeamPage(ctx echo.Context) error {
	return render(ctx, admin.Team())
}
func (a AdminHandler) PlansPage(ctx echo.Context) error {
	P := []models.Plan{{ID: 0, Name: "Admin", Description: "Route for administrators", Steps: []models.Step{{ID: 0, Description: "Create password", Done: true}, {ID: 1, Description: "Create account"}}}, {ID: 1, Name: "Sales", Description: "Route for salesman", Steps: []models.Step{{ID: 0, Description: "Create password"}, {ID: 1, Description: "Create account"}}}}
	return render(ctx, admin.Plans(P))
}
