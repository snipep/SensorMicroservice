package main

import (
	"database/sql"
	"net/http"
	"github.com/labstack/echo/v4"
)

type getSensorRequest struct {
	ID1 string `json:"id1"`
	ID2 int32  `json:"id2"`
}
func (app *Application) getSensorByIDs(c echo.Context) error {
	var req getSensorRequest
	if err := readJSON(c, &req); err != nil {
		app.logger.Errorw("failed to read or parse JSON", "error", err)
		return app.badRequestResponse(c, err)
	}


	sensor, err := app.store.Sensor.GetSensorByIDs(req.ID1, req.ID2)
	if err != nil {
		// Check for the specific "not found" error
		if err == sql.ErrNoRows {
			// Use "return" to stop execution here
			return app.notFoundResponse(c, err)
		}
		// For all other errors, it's an internal server error
		app.logger.Errorw("failed to fetch sensor", "id1", req.ID1, "id2", req.ID2, "error", err)
		return app.internalServerError(c, err)
	}

	err = app.jsonResponse(c, http.StatusOK, sensor)
	if err != nil {
		app.logger.Errorw("failed to write JSON response", "id1", req.ID1, "id2", req.ID2, "error", err)
		return app.internalServerError(c, err)
	}
	app.logger.Infow("sensor data fetched successfully", "id1", req.ID1, "id2", req.ID2)
	return nil
}