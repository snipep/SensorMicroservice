package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/snipep/Assessment/Consumer-Service/internal/store"
	"go.uber.org/zap"
)

type Application struct {
	config Config
	store *store.Storage
	logger *zap.SugaredLogger	
}


type dbConfig struct {
	addr string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime string
}

type Config struct {
	addr string
	db dbConfig
}


func (app *Application) RegiterRoutes(e *echo.Echo) {
	// --Routes--
	apiGroup := e.Group("api/v1")
	sensordata := apiGroup.Group("/sensordata")
	sensordata.GET("/query", app.getSensorByIDs)
	sensordata.GET("/history", app.getSensorHistory)
	sensordata.GET("/query/history", app.getSensorHistoryByIDs)

}

func (app *Application) run(echo *echo.Echo) error {
	app.RegiterRoutes(echo)
	srv := http.Server{
		Addr:    app.config.addr,
		Handler: echo,
		WriteTimeout: 15 * time.Second,
		ReadTimeout: 15 * time.Second,
		IdleTimeout: time.Minute,
	}

	app.logger.Infow("Server has started", "addr", app.config.addr)
	return srv.ListenAndServe()
}
