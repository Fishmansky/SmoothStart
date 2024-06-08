package handlers

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func render(c echo.Context, cmp templ.Component) error {
	return cmp.Render(c.Request().Context(), c.Response())
}
