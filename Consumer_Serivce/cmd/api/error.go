package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HTTPError struct {
    Error string `json:"error"`
}

func (app *Application) internalServerError(c echo.Context, err error) error {
	app.logger.Errorw("internal error",
		"method", c.Request().Method,
		"path", c.Request().URL.Path,
		"error", err.Error(),
	)
	return writeJSONError(c, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *Application) badRequestResponse(c echo.Context, err error) error {
	app.logger.Warnw("bad request",
		"method", c.Request().Method,
		"path", c.Request().URL.Path,
		"error", err.Error(),
	)
	return writeJSONError(c, http.StatusBadRequest, err.Error())
}

func (app *Application) notFoundResponse(c echo.Context, err error) error {
	app.logger.Errorw("not found",
		"method", c.Request().Method,
		"path", c.Request().URL.Path,
		"error", err.Error(),
	)
	return writeJSONError(c, http.StatusNotFound, "Record not found")
}

func (app *Application) conflictResponse(c echo.Context, err error) error {
	app.logger.Warnw("conflict response",
		"method", c.Request().Method,
		"path", c.Request().URL.Path,
		"error", err.Error(),
	)
	return writeJSONError(c, http.StatusConflict, "conflict occurred")
}
