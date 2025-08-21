package main

import (
	"context"
	"net/http"
	"time"
	_ "github.com/snipep/Assessment/Consumer-Service/docs" 
	"github.com/swaggo/echo-swagger"

	"github.com/labstack/echo/v4"
	"github.com/snipep/Assessment/Consumer-Service/internal/store"
	"go.uber.org/zap"
)

type Application struct {
	config Config
	store  *store.Storage
	logger *zap.SugaredLogger
	// worker pool
	ingestCh    chan store.SensorData
	workerCount int
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type Config struct {
	addr string
	db   dbConfig
}

func (app *Application) RegiterRoutes(e *echo.Echo) {

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// --Routes--
	apiGroup := e.Group("api/v1")
	// --- Auth Routes ---
	apiGroup.POST("/signup", app.signup)
	apiGroup.POST("/signin", app.signin)
	
	// Protected with JWT in main.go middleware
	sensordata := apiGroup.Group("/sensordata")
	// --- GET Routes ---
	sensordata.GET("/query", app.getSensorByIDs)
	sensordata.GET("/history", app.getSensorHistory)
	sensordata.GET("/query-history", app.getSensorHistoryByIDs)

	// --- DELETE Routes ---
	sensordata.DELETE("/query", app.deleteSensorDataByIDs)
	sensordata.DELETE("/history", app.deleteSensorHistory)
	sensordata.DELETE("/query-history", app.deleteSensorHistoryByIDs)

	// --- PUT Routes ---
	sensordata.PUT("/query", app.editSensorDataByID)
	sensordata.PUT("/history", app.editSensorHistory)
	sensordata.PUT("/query-history", app.editSensorHistoryByIDs)

}

func (app *Application) run(echo *echo.Echo) error {
	app.RegiterRoutes(echo)
	srv := http.Server{
		Addr:         app.config.addr,
		Handler:      echo,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  time.Minute,
	}

	app.logger.Infow("Server has started", "addr", app.config.addr)
	return srv.ListenAndServe()
}

// InitIngest initializes the buffered channel and starts worker goroutines.
func (app *Application) InitIngest(workers int, bufferSize int) {
	if workers <= 0 {
		workers = 4
	}
	if bufferSize <= 0 {
		bufferSize = 1024
	}
	app.workerCount = workers
	app.ingestCh = make(chan store.SensorData, bufferSize)
	for i := 0; i < app.workerCount; i++ {
		go app.ingestWorker(i)
	}
	app.logger.Infof("Ingest worker pool started with %d workers (buffer=%d)", app.workerCount, bufferSize)
}

func (app *Application) ingestWorker(id int) {
	for data := range app.ingestCh {
		// Use a short-lived context per insert
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		if err := app.store.Sensor.InsertSensorData(ctx, data); err != nil {
			app.logger.Errorf("worker %d: failed to insert sensor data: %v", id, err)
		}
		cancel()
	}
}
