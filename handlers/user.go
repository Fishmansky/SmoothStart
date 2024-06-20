package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
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
func (u UserHandler) getUserID(c echo.Context) (int, error) {
	var uID int
	userName, err := c.Cookie("username")
	if err != nil {
		return uID, err
	}
	if err := u.db.QueryRow("SELECT id  FROM users WHERE username = $1", userName.Value).Scan(&uID); err != nil {
		if err == sql.ErrNoRows {
			return uID, fmt.Errorf("%s\n", "User not found")
		}
		return uID, fmt.Errorf("%s\n", err)
	}
	return uID, nil
}
func (u UserHandler) getPlan(id int) (models.Plan, error) {
	var plan models.Plan
	var s []byte
	if err := u.db.QueryRow("Select id, name, description, steps FROM plans WHERE user_id = $1", id).Scan(&plan.ID, &plan.Name, &plan.Description, &s); err != nil {
		slog.Warn(err.Error())
		return plan, err
	}
	var steps []models.Step
	json.Unmarshal(s, &steps)
	plan.Steps = append(plan.Steps, steps...)
	return plan, nil
}

func (u UserHandler) UpdateStepStatus(userId int, stepId int) error {
	var steps []models.Step
	var s []byte
	if err := u.db.QueryRow("Select steps FROM plans WHERE user_id = $1", userId).Scan(&s); err != nil {
		return err
	}
	json.Unmarshal(s, &steps)
	if steps[stepId].Done {
		steps[stepId].Done = false
		updatedSteps, err := json.Marshal(&steps)
		if err != nil {
			return err
		}
		_, err = u.db.Exec("UPDATE plans SET steps = $1 WHERE user_id = $2", updatedSteps, userId)
		if err != nil {
			return err
		}
	} else {
		steps[stepId].Done = true
		updatedSteps, err := json.Marshal(&steps)
		if err != nil {
			return err
		}
		_, err = u.db.Exec("UPDATE plans SET steps = $1 WHERE user_id = $2", updatedSteps, userId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u UserHandler) HomePage(c echo.Context) error {
	return render(c, user.Home())
}

func (u UserHandler) PlanPage(c echo.Context) error {
	userId, err := u.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	p, err := u.getPlan(userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if isHtmx(c) {
		return render(c, user.Plan(p))
	}
	return render(c, user.PlanPage(p))
}

func (u UserHandler) HandleUpdateStepStatus(c echo.Context) error {
	sId := c.QueryParam("step")
	userId, err := u.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	stepId, err := strconv.Atoi(sId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	err = u.UpdateStepStatus(userId, stepId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "Status updated")

}
