package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (app *Application) RegiterRoutes(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// --Routes--
	apiGroup := e.Group("api/v1")

	apiGroup.POST("/frequency", echo.WrapHandler(http.HandlerFunc(app.UpdateFrequencyHandler)))
}
