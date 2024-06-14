package handlers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
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

func (a AdminHandler) HomePage(c echo.Context) error {
	var data []models.DashboardData
	rows, err := a.db.Query("SELECT id, fname, sname FROM users WHERE is_admin = '0'")
	if err != nil {
		return c.JSON(http.StatusNoContent, "No teammates found")
	}
	defer rows.Close()

	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.ID, &u.Fname, &u.Sname)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Error scanning user sql")
		}
		planRows, err := a.db.Query("SELECT id, name, steps FROM plans WHERE userid = $1", u.ID)
		if err != nil {
			return c.JSON(http.StatusNoContent, "No plans found")
		}
		defer planRows.Close()

		var p models.Plan
		var s []byte
		for planRows.Next() {
			err := planRows.Scan(&p.ID, &p.Name, &s)
			if err != nil {
				slog.Error(err.Error())
				return c.JSON(http.StatusInternalServerError, "Error scanning plan sql")
			}

		}
		var steps []models.Step
		json.Unmarshal(s, &steps)
		p.Steps = append(p.Steps, steps...)
		data = append(data, models.DashboardData{u, p.CompletionStatus()})
	}

	return render(c, admin.Home(data))
}
func (a AdminHandler) TeamPage(c echo.Context) error {
	return render(c, admin.Team())
}
func (a AdminHandler) PlansPage(c echo.Context) error {
	var plans []models.Plan
	rows, err := a.db.Query("Select * FROM plans")
	if err != nil {
		slog.Info(err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var p models.Plan
		var s []byte
		var uId int
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &uId, &s)
		if err != nil {
			slog.Error(err.Error())
			return c.JSON(http.StatusInternalServerError, "Error scanning plan sql")
		}
		var steps []models.Step
		json.Unmarshal(s, &steps)
		p.Steps = append(p.Steps, steps...)
		plans = append(plans, p)
	}
	return render(c, admin.Plans(plans))
}
