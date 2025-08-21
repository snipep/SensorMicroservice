package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// --- Auth Hanlders ---
type signupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type signupResponse struct {
    Message   string    `json:"message"`
    CreatedAt time.Time `json:"created_at"`
}

// signup godoc
// @Summary      User Signup
// @Description  Create a new user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      signupRequest  true  "User Signup Request"
// @Success      201    {object}  signupResponse
// @Failure      400    {object}  HTTPError
// @Failure      500    {object}  HTTPError
// @Router       /signup [post]
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
	resp := signupResponse{
		Message:    "user created",
		CreatedAt: time.Now().UTC(),
	}
	return app.jsonResponse(c, http.StatusCreated, resp)
}

type signinRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type signinResponse struct {
	Token string `json:"token"`
}

// signin godoc
// @Summary      User Signin
// @Description  Authenticate user and return JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      signinRequest  true  "User Signin Request"
// @Success      200    {object}  signinResponse
// @Failure      400    {object}  HTTPError
// @Failure      500    {object}  HTTPError
// @Router       /signin [post]
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
	resp := signinResponse{
		Token: token,
	}
	return app.jsonResponse(c, http.StatusOK, resp)
}


// getSensorByIDs godoc
// @Summary      Get sensor data by IDs
// @Description  Responds with sensor data for the given IDs. Protected by JWT.
// @Tags         sensordata - GET
// @Accept       json
// @Produce      json
// @Param        id1    query     string  false  "Sensor ID 1"
// @Param        id2    query     int32   false  "Sensor ID 2"
// @Param        limit  query     int     false  "Limit for pagination" default(10)
// @Param        offset query     int     false  "Offset for pagination" default(0)
// @Success      200    {object}  store.SensorData
// @Failure      400    {object}  HTTPError
// @Failure      401    {object}  HTTPError
// @Failure      404    {object}  HTTPError
// @Failure      500    {object}  HTTPError
// @Security     Bearer
// @Router       /sensordata/query [get]
func (app *Application) getSensorByIDs(c echo.Context) error {
	ID1, ID2, _ := app.getIDsByParams(c)
	limit, offset, _ := app.applyDefaultPagination(c)
	sensor, err := app.store.Sensor.GetSensorByIDs(ID1, int32(ID2), limit, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return app.notFoundResponse(c, err)
		}
		return app.internalServerError(c, err)
	}

	return app.jsonResponse(c, http.StatusOK, sensor)
}

// getSensorHistory godoc
// @Summary      Get sensor history
// @Description  Responds with sensor history data within the specified time range. Protected by JWT.
// @Tags         sensordata - GET
// @Accept       json
// @Produce      json
// @Param        start_time    query     string  true  "Start time in RFC3339 format"
// @Param        end_time      query     string  true  "End time in RFC3339 format"
// @Param        limit         query     int     false  "Limit for pagination" default(10)
// @Param        offset        query     int     false  "Offset for pagination" default(0)
// @Success      200    {object}  []store.SensorData
// @Failure      400    {object}  HTTPError
// @Failure      401    {object}  HTTPError
// @Failure      404    {object}  HTTPError
// @Failure      500    {object}  HTTPError
// @Security     Bearer
// @Router       /sensordata/history [get]
func (app *Application) getSensorHistory(c echo.Context) error {
	StartTime, EndTime, _ := app.getTimestampByParams(c)
	limit, offset, _ := app.applyDefaultPagination(c)

	// Parse the start and end time strings into time.Time objects
	startTime, err := time.Parse(time.RFC3339, StartTime)
	if err != nil {
		return app.badRequestResponse(c, err)
	}

	endTime, err := time.Parse(time.RFC3339, EndTime)
	if err != nil {
		return app.badRequestResponse(c, err)
	}

	// Call the new store method
	sensors, err := app.store.Sensor.GetSensorHistory(startTime, endTime, limit, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			app.notFoundResponse(c, err)
			return err
		}
		return app.internalServerError(c, err)
	}

	return app.jsonResponse(c, http.StatusOK, sensors)
}


// getSensorHistoryByIDs godoc
// @Summary      Get sensor history by IDs
// @Description  Responds with sensor history data for the specified IDs within a time range. Protected by JWT.
// @Tags         sensordata - GET
// @Accept       json
// @Produce      json
// @Param        id1         query     string  true  "Sensor ID 1"
// @Param        id2         query     int32   true  "Sensor ID 2"
// @Param        start_time  query     string  true  "Start time in RFC3339 format"
// @Param        end_time    query     string  true  "End time in RFC3339 format"
// @Param        limit       query     int     false  "Limit for pagination" default(10)
// @Param        offset      query     int     false  "Offset for pagination" default(0)
// @Success      200      {object}  []store.SensorData
// @Failure      400      {object}  HTTPError
// @Failure      401      {object}  HTTPError
// @Failure      404      {object}  HTTPError
// @Failure      500      {object}  HTTPError
// @Security     Bearer
// @Router       /sensordata/query-history [get]
func (app *Application) getSensorHistoryByIDs(c echo.Context) error {
	ID1, ID2, _ := app.getIDsByParams(c)
	StartTime, EndTime, _ := app.getTimestampByParams(c)
	limit, offset, _ := app.applyDefaultPagination(c)

	// Parse the start and end time strings
	startTime, err := time.Parse(time.RFC3339, StartTime)
	if err != nil {
		return app.badRequestResponse(c, err)
	}

	endTime, err := time.Parse(time.RFC3339, EndTime)
	if err != nil {
		return app.badRequestResponse(c, err)
	}

	sensors, err := app.store.Sensor.GetSensorHistoryByIDs(ID1, int32(ID2), startTime, endTime, limit, offset)
	log.Println("sensors:", sensors)
	if err != nil {
		return app.internalServerError(c, err)
	}

	return app.jsonResponse(c, http.StatusOK, sensors)
}


type SensorRequest struct {
	ID1    string `json:"id1"`
	ID2    int32  `json:"id2"`
}

type rowsAffectedResponse struct {
	Message      string `json:"message"`
	RowsAffected int64  `json:"rows_affected"`
}

// deleteSensorDataByIDs godoc
// @Summary      Delete sensor data by IDs
// @Description  Deletes sensor data records matching the provided ID1 and ID2. Protected by JWT.
// @Tags         sensordata - DELETE
// @Accept       json
// @Produce      json
// @Param        request  body      SensorRequest  true  "Sensor IDs to delete"
// @Success      200      {object}  rowsAffectedResponse
// @Failure      400      {object}  HTTPError
// @Failure      401      {object}  HTTPError
// @Failure      500      {object}  HTTPError
// @Security     Bearer
// @Router       /sensordata/query [delete]
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

	resp := rowsAffectedResponse{
		Message:       "Sensor data deleted successfully",
		RowsAffected: rowsAffected,
	}
	return app.jsonResponse(c, http.StatusOK, resp)
}

type SensorHistoryRequest struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// deleteSensorHistory godoc	
// @Summary      Delete sensor history
// @Description  Deletes sensor history records within the specified time range. Protected by JWT.
// @Tags         sensordata - DELETE
// @Accept       json
// @Produce      json
// @Param        request  body      SensorHistoryRequest  true  "Time range for sensor history deletion"
// @Success      200      {object}  rowsAffectedResponse
// @Failure      400      {object}  HTTPError
// @Failure      401      {object}  HTTPError
// @Failure      500      {object}  HTTPError
// @Security     Bearer
// @Router       /sensordata/history [delete]
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

	resp := rowsAffectedResponse{
		Message:      "Sensor history deleted successfully",
		RowsAffected: rowsAffected,
	}
	return app.jsonResponse(c, http.StatusOK, resp)
}

type SensorHistoryByIDsRequest struct {
	ID1       string `json:"id1"`
	ID2       int32  `json:"id2"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// deleteSensorHistoryByIDs godoc
// @Summary      Delete sensor history by IDs
// @Description  Deletes sensor history records for the specified IDs within a time range. Protected by JWT.
// @Tags         sensordata - DELETE
// @Accept       json
// @Produce      json
// @Param        request  body      SensorHistoryByIDsRequest  true  "Sensor IDs and time range for deletion"
// @Success      200      {object}  rowsAffectedResponse
// @Failure      400      {object}  HTTPError
// @Failure      401      {object}  HTTPError
// @Failure      500      {object}  HTTPError
// @Security     Bearer
// @Router       /sensordata/query-history [delete]
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

	resp := rowsAffectedResponse{
		Message:       "Sensor history for specified IDs deleted successfully",
		RowsAffected: rowsAffected,
	}
	return app.jsonResponse(c, http.StatusOK, resp)
}

type editByIDRequest struct {
	ID1      string  `json:"id1"`
	ID2      int32   `json:"id2"`
	NewValue float32 `json:"new_value"`
}

// editSensorDataByID godoc
// @Summary      Edit sensor data by ID
// @Description  Updates the value of a specific sensor data record. Protected by JWT.
// @Tags         sensordata - PUT
// @Accept       json
// @Produce      json
// @Param        request  body      editByIDRequest  true  "Sensor ID and new value"
// @Success      200      {object}  rowsAffectedResponse
// @Failure      400      {object}  HTTPError
// @Failure      401      {object}  HTTPError
// @Failure      500      {object}  HTTPError
// @Security     Bearer
// @Router       /sensordata/query [put]
func (app *Application) editSensorDataByID(c echo.Context) error {
	var req editByIDRequest
	if err := readJSON(c, &req); err != nil {
		return app.badRequestResponse(c, err)
	}

	rowsAffected, err := app.store.Sensor.EditSensorDataByID(req.ID1, req.ID2, req.NewValue)
	if err != nil {
		return app.internalServerError(c, err)
	}

	resp := rowsAffectedResponse{
		Message:       "Sensor data updated successfully",
		RowsAffected: rowsAffected,
	}
	return app.jsonResponse(c, http.StatusOK, resp)
}

type editByHistoryRequest struct {
	StartTime string  `json:"start_time"`
	EndTime   string  `json:"end_time"`
	NewValue  float32 `json:"new_value"`
}

// editSensorHistory godoc
// @Summary      Edit sensor history
// @Description  Updates sensor history records within a specified time range. Protected by JWT.
// @Tags         sensordata - PUT
// @Accept       json
// @Produce      json
// @Param        request  body      editByHistoryRequest  true  "Time range and new value for sensor history"
// @Success      200      {object}  rowsAffectedResponse
// @Failure      400      {object}  HTTPError
// @Failure      401      {object}  HTTPError
// @Failure      500      {object}  HTTPError
// @Security     Bearer
// @Router       /sensordata/history [put]
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

	resp := rowsAffectedResponse{
		Message:      "Sensor history updated successfully",
		RowsAffected: rowsAffected,
	}
	return app.jsonResponse(c, http.StatusOK, resp)
}

type editByQueryHistoryRequest struct {
	ID1       string  `json:"id1"`
	ID2       int32   `json:"id2"`
	StartTime string  `json:"start_time"`
	EndTime   string  `json:"end_time"`
	NewValue  float32 `json:"new_value"`
}

// editSensorHistoryByIDs godoc
// @Summary      Edit sensor history by IDs
// @Description  Updates sensor history records for specified IDs within a time range. Protected by JWT.
// @Tags         sensordata - PUT
// @Accept       json
// @Produce      json
// @Param        request  body      editByQueryHistoryRequest  true  "Sensor IDs and time range for update"
// @Success      200      {object}  rowsAffectedResponse
// @Failure      400      {object}  HTTPError
// @Failure      401      {object}  HTTPError
// @Failure      500      {object}  HTTPError
// @Security     Bearer
// @Router       /sensordata/query-history [put]
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

	resp := rowsAffectedResponse{
		Message:       "Sensor history for specified IDs updated successfully",
		RowsAffected: rowsAffected,
	}
	return app.jsonResponse(c, http.StatusOK, resp)
}

func (app *Application) applyDefaultPagination(c echo.Context) (limit, offset int, err error) {
	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")

	// If a limit was provided, try to parse it.
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return 0, 0,app.badRequestResponse(c, fmt.Errorf("invalid limit parameter: must be a number"))
		}
	}

	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			return 0, 0, app.badRequestResponse(c, fmt.Errorf("invalid offset parameter: must be a number"))
		}
	}
	if limit <= 0 {
		limit = 10 
	}
	if offset < 0 {
		offset = 0 
	}
	return limit, offset, nil
}


func (app *Application) getIDsByParams(c echo.Context) (string, int, error) {
	ID1 := c.QueryParam("id1")
	ID2, err := strconv.Atoi(c.QueryParam("id2"))
	if err != nil {
		return "", 0, app.badRequestResponse(c, fmt.Errorf("invalid id2 parameter: must be a number"))
	}
	return ID1, ID2, nil
}

func (app *Application) getTimestampByParams(c echo.Context) (string, string, error) {
	startTime := c.QueryParam("start_time")
	endTime := c.QueryParam("end_time")
	return startTime, endTime, nil
}
