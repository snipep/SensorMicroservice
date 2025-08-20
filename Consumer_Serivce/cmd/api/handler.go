package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

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


// --- NEW: Define a struct for the history request payload ---
type getSensorHistoryRequest struct {
	// Using RFC3339 format for timestamps, e.g., "2020-05-04T00:00:00Z"
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// --- NEW: Handler for fetching data by time duration ---
func (app *Application) getSensorHistory(c echo.Context) error {
	var req getSensorHistoryRequest
	if err := readJSON(c, &req); err != nil {
		app.logger.Errorw("failed to read or parse JSON for history request", "error", err)
		return app.badRequestResponse(c, err)
	}

	// Parse the start and end time strings into time.Time objects
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		app.logger.Errorw("invalid start_time format", "start_time", req.StartTime, "error", err)
		return app.badRequestResponse(c, err)
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		app.logger.Errorw("invalid end_time format", "end_time", req.EndTime, "error", err)
		return app.badRequestResponse(c, err)
	}

	// Call the new store method
	sensors, err := app.store.Sensor.GetSensorHistory(startTime, endTime)
	log.Println("sensors:", sensors)
	if err != nil {
		if err == sql.ErrNoRows {
			app.notFoundResponse(c, err)
			return err
		}
		// This will handle cases where no rows are found gracefully (returns an empty list)
		app.logger.Errorw("failed to fetch sensor history", "start", req.StartTime, "end", req.EndTime, "error", err)
		return app.internalServerError(c, err)
	}

	// Write the successful response
	err = app.jsonResponse(c, http.StatusOK, sensors)
	if err != nil {
		app.logger.Errorw("failed to write JSON response for history", "error", err)
		return app.internalServerError(c, err)
	}

	app.logger.Infow("sensor history fetched successfully", "start", req.StartTime, "end", req.EndTime)
	return nil
}
