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
	SensorClient sensor.SensorServiceClient
	logger 	 	*zap.SugaredLogger
	ticker       *time.Ticker
	tickerMutex  *sync.Mutex
}

func NewApplication(sensorClient sensor.SensorServiceClient, ticker *time.Ticker, logger *zap.SugaredLogger) *Application {
	return &Application{
		SensorClient: sensorClient,
		ticker:       ticker,
		logger: logger,
		tickerMutex:  &sync.Mutex{},
	}
}

type FrequencyUpdateRequest struct {
	FrequencySeconds int `json:"frequency_seconds"`
}

// UpdateFrequencyHandler handles the frequency update request.
func (app *Application) UpdateFrequencyHandler(w http.ResponseWriter, r *http.Request){
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
	app.ticker.Reset(time.Duration(req.FrequencySeconds) * time.Second) // Reset the ticker to the new frequency.
	// app.ticker.Stop() // Stop the old ticker first.
	// app.ticker = time.NewTicker(time.Duration(req.FrequencySeconds) * time.Second) // Replace it with a new one.
	app.tickerMutex.Unlock()

	app.logger.Infof("Frequency updated to %d seconds", req.FrequencySeconds)
	writeJSON(w, http.StatusOK, map[string]string{"message": "Frequency updated successfully"})
}