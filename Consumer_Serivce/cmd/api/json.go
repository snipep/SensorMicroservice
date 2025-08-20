package main

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

func writeJSON(c echo.Context, status int, data any) error {
	
	return c.JSON(status, data)
}

func readJSON(c echo.Context, data any) error {
	// Limit body size (1MB max)
	c.Request().Body = http.MaxBytesReader(c.Response(), c.Request().Body, 1_048_576)

	decoder := json.NewDecoder(c.Request().Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

func writeJSONError(c echo.Context, status int, message string) error {
	type envelope struct {
		Error string `json:"error"`
	}
	return writeJSON(c, status, &envelope{Error: message})
}

func (app *Application) jsonResponse(c echo.Context, status int, data any) error {
	type envelop struct {
		Data any `json:"data"`
	}
	return writeJSON(c, status, &envelop{Data: data})
}