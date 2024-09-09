package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
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
	result := c.Get("username")
	uName, ok := result.(string)
	if ok == false {
		return -1, errors.New("Unable to retrieve username")
	}
	var user models.User
	if err := u.db.QueryRow("SELECT id, username, fname, lname, is_admin FROM users WHERE username = $1", uName).Scan(&user.ID, &user.Username, &user.Fname, &user.Lname, &user.IsAdmin); err != nil {
		if err == sql.ErrNoRows {
			return -1, fmt.Errorf("%s\n", "User not found")
		}
		return -1, fmt.Errorf("%s\n", err)
	}
	return user.ID, nil
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

func (u UserHandler) UpdateStepStatus(userId int, stepID int) (*models.Step, error) {
	var steps []models.Step
	var s []byte
	if err := u.db.QueryRow("Select steps FROM plans WHERE user_id = $1", userId).Scan(&s); err != nil {
		return nil, err
	}
	json.Unmarshal(s, &steps)
	if steps[stepID].Done {
		steps[stepID].Done = false
		updatedSteps, err := json.Marshal(&steps)
		if err != nil {
			return nil, err
		}
		_, err = u.db.Exec("UPDATE plans SET steps = $1 WHERE user_id = $2", updatedSteps, userId)
		if err != nil {
			return nil, err
		}
	} else {
		steps[stepID].Done = true
		updatedSteps, err := json.Marshal(&steps)
		if err != nil {
			return nil, err
		}
		_, err = u.db.Exec("UPDATE plans SET steps = $1 WHERE user_id = $2", updatedSteps, userId)
		if err != nil {
			return nil, err
		}
	}
	return &steps[stepID], nil
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
	return render(c, user.PlanPage(p))
}

func (u UserHandler) HandleUpdateStepStatus(c echo.Context) error {
	sID := c.FormValue("step")
	pID := c.FormValue("plan")
	userID, err := u.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	fmt.Println(sID)
	stepID, err := strconv.Atoi(sID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	planID, err := strconv.Atoi(pID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	step, err := u.UpdateStepStatus(userID, stepID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return render(c, user.PlanStep(step.ID, planID, step.Done, step.Description))

}
