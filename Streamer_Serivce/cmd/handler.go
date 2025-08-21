package main

import (
	"errors"
	"net/http"
	"sync"
	"time"

	sensor "github.com/snipep/Assessment/Streamer-Service/Proto"
	"go.uber.org/zap"
)

type Application struct {
	SensorClient    sensor.SensorServiceClient
	logger          *zap.SugaredLogger
	ticker          *time.Ticker
	tickerMutex     *sync.Mutex
	FixedSensorID1  string
	FixedSensorType string
}

func NewApplication(sensorClient sensor.SensorServiceClient, ticker *time.Ticker, logger *zap.SugaredLogger, sensorID1 string, sensorType string) *Application {
	return &Application{
		SensorClient:    sensorClient,
		ticker:          ticker,
		logger:          logger,
		tickerMutex:     &sync.Mutex{},
		FixedSensorID1:  sensorID1,
		FixedSensorType: sensorType,
	}
}

type FrequencyUpdateRequest struct {
	FrequencySeconds int `json:"frequency_seconds"`
}

func (app *Application) UpdateFrequencyHandler(w http.ResponseWriter, r *http.Request) {
	var req FrequencyUpdateRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if req.FrequencySeconds <= 0 {
		app.badRequestResponse(w, r, errors.New("frequency must be greater than zero"))
		return
	}

	// Safely update the ticker
	app.tickerMutex.Lock()
	app.ticker.Reset(time.Duration(req.FrequencySeconds) * time.Second)
	app.tickerMutex.Unlock()

	app.logger.Infof("Frequency updated to %d seconds", req.FrequencySeconds)
	writeJSON(w, http.StatusOK, map[string]string{"message": "Frequency updated successfully"})
}
