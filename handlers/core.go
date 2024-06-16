package handlers

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func render(c echo.Context, cmp templ.Component) error {
	return cmp.Render(c.Request().Context(), c.Response())
}

func isHtmx(c echo.Context) bool {
	return c.Request().Header.Get("HX-Request") == "true"
}
