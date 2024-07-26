package handlers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"smoothstart/models"
	"smoothstart/views/admin"
	"strconv"
	"strings"

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

func (a AdminHandler) GetTemplatePlans() ([]models.Plan, error) {
	var plans []models.Plan
	rows, err := a.db.Query("Select * FROM plan_templates")
	if err != nil {
		return plans, err
	}
	defer rows.Close()
	for rows.Next() {
		var p models.Plan
		var s []byte
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &s)
		if err != nil {
			return plans, err
		}
		var steps []models.Step
		json.Unmarshal(s, &steps)
		p.Steps = append(p.Steps, steps...)
		plans = append(plans, p)
	}
	return plans, nil

}
func (a AdminHandler) GetTemplatePlan(i int) (models.Plan, error) {
	var plan models.Plan
	rows, err := a.db.Query("Select * FROM plan_templates WHERE id = $1", i)
	if err != nil {
		slog.Info(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var s []byte
		err := rows.Scan(&plan.ID, &plan.Name, &plan.Description, &s)
		if err != nil {
			return plan, err
		}
		var steps []models.Step
		json.Unmarshal(s, &steps)
		plan.Steps = append(plan.Steps, steps...)
	}
	return plan, nil
}

func (a AdminHandler) getMemberPlan(id int) (models.Plan, error) {
	var plan models.Plan
	var s []byte
	if err := a.db.QueryRow("Select id, name, description, steps FROM plans WHERE user_id = $1", id).Scan(&plan.ID, &plan.Name, &plan.Description, &s); err != nil {
		slog.Warn(err.Error())
		return plan, err
	}
	var steps []models.Step
	json.Unmarshal(s, &steps)
	plan.Steps = append(plan.Steps, steps...)
	return plan, nil
}

func (a AdminHandler) AddPlanTemplate(name string, desc string) error {
	s, err := json.Marshal([]models.Step{})
	if err != nil {
		return err
	}
	_, err = a.db.Exec("INSERT INTO plan_templates (name, description,steps) VALUES ($1, $2, $3)", name, desc, s)
	if err != nil {
		return err
	}
	return nil
}

func (a AdminHandler) EditPlanTemplate(id int, name string, desc string) error {
	_, err := a.db.Exec("UPDATE plan_templates SET name = $1, description = $2 WHERE id = $3", name, desc, id)
	if err != nil {
		return err
	}
	return nil
}

func (a AdminHandler) UpdateTemplateStep(planId int, stepId int, desc string) error {
	var steps []models.Step
	var s []byte
	if err := a.db.QueryRow("Select steps FROM plan_templates WHERE id = $1", planId).Scan(&s); err != nil {
		return err
	}
	json.Unmarshal(s, &steps)
	steps[stepId].Description = desc
	updatedSteps, err := json.Marshal(&steps)
	if err != nil {
		return err
	}
	_, err = a.db.Exec("UPDATE plan_templates SET steps = $1 WHERE id = $2", updatedSteps, planId)
	if err != nil {
		return err
	}
	return nil
}

func (a AdminHandler) UpdateStepStatus(planId int, stepId int) error {
	var steps []models.Step
	var s []byte
	if err := a.db.QueryRow("Select steps FROM plans WHERE id = $1", planId).Scan(&s); err != nil {
		return err
	}
	json.Unmarshal(s, &steps)
	if steps[stepId].Done {
		steps[stepId].Done = false
		updatedSteps, err := json.Marshal(&steps)
		if err != nil {
			return err
		}
		_, err = a.db.Exec("UPDATE plans SET steps = $1 WHERE id = $2", updatedSteps, planId)
		if err != nil {
			return err
		}
	} else {
		steps[stepId].Done = true
		updatedSteps, err := json.Marshal(&steps)
		if err != nil {
			return err
		}
		_, err = a.db.Exec("UPDATE plans SET steps = $1 WHERE id = $2", updatedSteps, planId)
		if err != nil {
			return err
		}
	}
	return nil
}
func (a AdminHandler) AddStepToPlan(id int, step string) (int, error) {
	var steps []models.Step
	var s []byte
	if err := a.db.QueryRow("Select steps FROM plan_templates WHERE id = $1", id).Scan(&s); err != nil {
		return -1, err
	}
	json.Unmarshal(s, &steps)
	var last int
	if len(steps) > 0 {
		last = steps[len(steps)-1].ID + 1
	} else {
		last = 0
	}
	steps = append(steps, models.Step{
		ID:          last,
		Description: step,
		Done:        false,
	})
	updatedSteps, err := json.Marshal(&steps)
	if err != nil {
		return -1, err
	}
	_, err = a.db.Exec("UPDATE plan_templates SET steps = $1 WHERE id = $2", updatedSteps, id)
	if err != nil {
		return -1, err
	}
	return last, nil
}

func (a AdminHandler) HomePage(c echo.Context) error {
	var data []admin.DashboardUserData
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
		up, err := a.getMemberPlan(u.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Error scanning user sql")
		}
		data = append(data, admin.DashboardUserData{u, up.CompletionStatus()})
	}

	return render(c, admin.Home(data))
}

func (a AdminHandler) TeamPage(c echo.Context) error {
	var users []models.User
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
		_, err = a.getMemberPlan(u.ID)
		if err != nil {
			u.HasPlan = false
			users = append(users, u)
		} else {
			u.HasPlan = true
			users = append(users, u)
		}
	}
	return render(c, admin.Team(users))
}
func (a AdminHandler) PlansPage(c echo.Context) error {
	data := admin.NewPlanPageViewModel()
	plans, err := a.GetTemplatePlans()
	if err != nil {
		slog.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	data.Plans = plans
	return render(c, admin.PlansPage(*data))
}

func (a AdminHandler) PlanPage(c echo.Context) error {
	id := c.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	plan, err := a.GetTemplatePlan(i)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if isHtmx(c) {
		return render(c, admin.Plan(plan))
	}
	return render(c, admin.PlanPage(plan))
}

func (a AdminHandler) AddPlanPage(c echo.Context) error {
	return render(c, admin.AddPlanPage())
}

func (a AdminHandler) AddPlan(c echo.Context) error {
	name := c.FormValue("name")
	desc := c.FormValue("desc")
	err := a.AddPlanTemplate(name, desc)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	c.Response().Header().Set("HX-Location", "/user/home")
	return c.Redirect(http.StatusMovedPermanently, "/admin/plans")
}

func (a AdminHandler) MemberPlanPage(c echo.Context) error {
	id := c.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	plan, err := a.getMemberPlan(i)
	if err != nil {
		slog.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, "Error scanning plan sql")
	}
	if isHtmx(c) {
		return render(c, admin.Plan(plan))
	}
	//if c.Request().Header.Get("HX-Request") != "true" {
	//	fmt.Println("is not htmx")
	//	return c.Redirect(http.StatusMovedPermanently, "/admin/plans")
	//}
	return render(c, admin.PlanPage(plan))
}

func (a AdminHandler) MemberPlanStepStatusUpdate(c echo.Context) error {
	id := c.Param("id")
	sId := c.QueryParam("step")
	planId, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	stepId, err := strconv.Atoi(sId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	err = a.UpdateStepStatus(planId, stepId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "Status updated")
}

func (a AdminHandler) AddTemplateStep(c echo.Context) error {
	desc := c.FormValue("description")
	cur := c.Request().Header.Get("HX-Current-URL")
	if cur == "" {
		return c.JSON(http.StatusBadRequest, "Bad request")
	}
	parts := strings.Split(cur, "/")
	i := parts[len(parts)-1]
	pid, err := strconv.Atoi(i)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	sid, err := a.AddStepToPlan(pid, desc)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return render(c, admin.TemplateStep(sid, pid, desc))
}

func (a AdminHandler) GetEditStep(c echo.Context) error {
	desc := c.FormValue("description")
	p := c.FormValue("plan")
	pid, err := strconv.Atoi(p)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	s := c.FormValue("step")
	sid, err := strconv.Atoi(s)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return render(c, admin.TemplateStepEdit(sid, pid, desc))
}

func (a AdminHandler) EditStep(c echo.Context) error {
	desc := c.FormValue("description")
	p := c.FormValue("plan")
	pid, err := strconv.Atoi(p)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	s := c.FormValue("step")
	sid, err := strconv.Atoi(s)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	err = a.UpdateTemplateStep(pid, sid, desc)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return render(c, admin.TemplateStep(sid, pid, desc))
}

func (a AdminHandler) EditTemplate(c echo.Context) error {
	i := c.FormValue("id")
	id, err := strconv.Atoi(i)
	if err != nil {
		return err
	}
	name := c.FormValue("name")
	desc := c.FormValue("desc")
	err = a.EditPlanTemplate(id, name, desc)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "Plan modified")
}
