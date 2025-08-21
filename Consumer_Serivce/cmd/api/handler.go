package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type SensorRequest struct {
	ID1    string `json:"id1"`
	ID2    int32  `json:"id2"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

type SensorHistoryRequest struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
}

type SensorHistoryByIDsRequest struct {
	ID1       string `json:"id1"`
	ID2       int32  `json:"id2"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
}

func (app *Application) getSensorByIDs(c echo.Context) error {
	var req SensorRequest
	if err := readJSON(c, &req); err != nil {
		return app.badRequestResponse(c, err)
	}
	applyDefaultPagination(&req.Limit, &req.Offset)

	sensor, err := app.store.Sensor.GetSensorByIDs(req.ID1, req.ID2, req.Limit, req.Offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return app.notFoundResponse(c, err)
		}
		return app.internalServerError(c, err)
	}

	return app.jsonResponse(c, http.StatusOK, sensor)
}

func (app *Application) getSensorHistory(c echo.Context) error {
	var req SensorHistoryRequest
	if err := readJSON(c, &req); err != nil {
		return app.badRequestResponse(c, err)
	}
	applyDefaultPagination(&req.Limit, &req.Offset)

	// Parse the start and end time strings into time.Time objects
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return app.badRequestResponse(c, err)
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		return app.badRequestResponse(c, err)
	}

	// Call the new store method
	sensors, err := app.store.Sensor.GetSensorHistory(startTime, endTime, req.Limit, req.Offset)
	if err != nil {
		if err == sql.ErrNoRows {
			app.notFoundResponse(c, err)
			return err
		}
		return app.internalServerError(c, err)
	}

	return app.jsonResponse(c, http.StatusOK, sensors)
}

func (app *Application) getSensorHistoryByIDs(c echo.Context) error {
	var req SensorHistoryByIDsRequest
	if err := readJSON(c, &req); err != nil {
		return app.badRequestResponse(c, err)
	}
	applyDefaultPagination(&req.Limit, &req.Offset)

	// Parse the start and end time strings
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return app.badRequestResponse(c, err)
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		return app.badRequestResponse(c, err)
	}

	sensors, err := app.store.Sensor.GetSensorHistoryByIDs(req.ID1, req.ID2, startTime, endTime, req.Limit, req.Offset)
	log.Println("sensors:", sensors)
	if err != nil {
		return app.internalServerError(c, err)
	}

	return app.jsonResponse(c, http.StatusOK, sensors)
}

func (app *Application) deleteSensorDataByIDs(c echo.Context) error {
	var req SensorRequest
	if err := readJSON(c, &req); err != nil {
		return app.badRequestResponse(c, err)
	}

	if req.ID2 < 0 || req.ID1 == "" {
		return app.badRequestResponse(c, fmt.Errorf("the 'ids' array cannot be empty"))
	}

	rowsAffected, err := app.store.Sensor.DeleteSensorDataByIDs(req.ID1, req.ID2)
	if err != nil {
		return app.internalServerError(c, err)
	}

	return app.jsonResponse(c, http.StatusOK, map[string]interface{}{
		"message":       "Sensor data deleted successfully",
		"rows_affected": rowsAffected,
	})
}

func (app *Application) deleteSensorHistory(c echo.Context) error {
	var req SensorHistoryRequest
	if err := readJSON(c, &req); err != nil {
		return app.badRequestResponse(c, err)
	}

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return app.badRequestResponse(c, fmt.Errorf("invalid start_time format: %w", err))
	}
	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		return app.badRequestResponse(c, fmt.Errorf("invalid end_time format: %w", err))
	}

	rowsAffected, err := app.store.Sensor.DeleteSensorHistory(startTime, endTime)
	if err != nil {
		return app.internalServerError(c, err)
	}

	return app.jsonResponse(c, http.StatusOK, map[string]interface{}{
		"message":       "Sensor history deleted successfully",
		"rows_affected": rowsAffected,
	})
}

func (app *Application) deleteSensorHistoryByIDs(c echo.Context) error {
	var req SensorHistoryByIDsRequest
	if err := readJSON(c, &req); err != nil {
		return app.badRequestResponse(c, err)
	}

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return app.badRequestResponse(c, fmt.Errorf("invalid start_time format: %w", err))
	}
	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		return app.badRequestResponse(c, fmt.Errorf("invalid end_time format: %w", err))
	}

	rowsAffected, err := app.store.Sensor.DeleteSensorHistoryByIDs(req.ID1, req.ID2, startTime, endTime)
	if err != nil {
		return app.internalServerError(c, err)
	}

	return app.jsonResponse(c, http.StatusOK, map[string]interface{}{
		"message":       "Sensor history for specified IDs deleted successfully",
		"rows_affected": rowsAffected,
	})
}

type editByIDRequest struct {
	ID1      string  `json:"id1"`
	ID2      int32   `json:"id2"`
	NewValue float32 `json:"new_value"`
}

type editByHistoryRequest struct {
	StartTime string  `json:"start_time"`
	EndTime   string  `json:"end_time"`
	NewValue  float32 `json:"new_value"`
}

type editByQueryHistoryRequest struct {
	ID1       string  `json:"id1"`
	ID2       int32   `json:"id2"`
	StartTime string  `json:"start_time"`
	EndTime   string  `json:"end_time"`
	NewValue  float32 `json:"new_value"`
}

func (app *Application) editSensorDataByID(c echo.Context) error {
	var req editByIDRequest
	if err := readJSON(c, &req); err != nil {
		return app.badRequestResponse(c, err)
	}

	rowsAffected, err := app.store.Sensor.EditSensorDataByID(req.ID1, req.ID2, req.NewValue)
	if err != nil {
		return app.internalServerError(c, err)
	}

	return app.jsonResponse(c, http.StatusOK, map[string]interface{}{
		"message":       "Sensor data updated successfully",
		"rows_affected": rowsAffected,
	})
}

func (app *Application) editSensorHistory(c echo.Context) error {
	var req editByHistoryRequest
	if err := readJSON(c, &req); err != nil {
		return app.badRequestResponse(c, err)
	}

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return app.badRequestResponse(c, fmt.Errorf("invalid start_time format: %w", err))
	}
	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		return app.badRequestResponse(c, fmt.Errorf("invalid end_time format: %w", err))
	}

	rowsAffected, err := app.store.Sensor.EditSensorHistory(startTime, endTime, req.NewValue)
	if err != nil {
		return app.internalServerError(c, err)
	}

	return app.jsonResponse(c, http.StatusOK, map[string]interface{}{
		"message":       "Sensor history updated successfully",
		"rows_affected": rowsAffected,
	})
}

func (app *Application) editSensorHistoryByIDs(c echo.Context) error {
	var req editByQueryHistoryRequest
	if err := readJSON(c, &req); err != nil {
		return app.badRequestResponse(c, err)
	}

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return app.badRequestResponse(c, fmt.Errorf("invalid start_time format: %w", err))
	}
	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		return app.badRequestResponse(c, fmt.Errorf("invalid end_time format: %w", err))
	}

	rowsAffected, err := app.store.Sensor.EditSensorHistoryByIDs(req.ID1, req.ID2, startTime, endTime, req.NewValue)
	if err != nil {
		return app.internalServerError(c, err)
	}

	return app.jsonResponse(c, http.StatusOK, map[string]interface{}{
		"message":       "Sensor history for specified IDs updated successfully",
		"rows_affected": rowsAffected,
	})
}
func applyDefaultPagination(limit, offset *int) {
	if *limit <= 0 {
		*limit = 10 // Default page size
	}
	if *offset < 0 {
		*offset = 0 // Default offset
	}
}

// --- Auth additions ---
// We keep imports minimal here and rely on main.go for JWT middleware.

type signupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type signinRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Application) signup(c echo.Context) error {
	var req signupRequest
	if err := readJSON(c, &req); err != nil {
		return app.badRequestResponse(c, err)
	}
	email := strings.TrimSpace(strings.ToLower(req.Email))
	if email == "" || req.Password == "" || strings.TrimSpace(req.Name) == "" {
		return app.badRequestResponse(c, fmt.Errorf("name, email and password are required"))
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return app.internalServerError(c, err)
	}
	_, err = app.store.User.CreateUser(c.Request().Context(), req.Name, email, string(hash))
	if err != nil {
		return app.internalServerError(c, err)
	}
	return app.jsonResponse(c, http.StatusCreated, map[string]any{
		"message":    "user created",
		"created_at": time.Now().UTC(),
	})
}

func (app *Application) signin(c echo.Context) error {
	var req signinRequest
	if err := readJSON(c, &req); err != nil {
		return app.badRequestResponse(c, err)
	}
	email := strings.TrimSpace(strings.ToLower(req.Email))
	if email == "" || req.Password == "" {
		return app.badRequestResponse(c, fmt.Errorf("email and password are required"))
	}
	user, err := app.store.User.GetUserByEmail(c.Request().Context(), email)
	if err != nil {
		if err == sql.ErrNoRows {
			return app.badRequestResponse(c, fmt.Errorf("invalid credentials"))
		}
		return app.internalServerError(c, err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return app.badRequestResponse(c, fmt.Errorf("invalid credentials"))
	}
	token, err := app.generateJWT(user.ID, user.Email)
	if err != nil {
		return app.internalServerError(c, err)
	}
	return app.jsonResponse(c, http.StatusOK, map[string]any{"token": token})
}
