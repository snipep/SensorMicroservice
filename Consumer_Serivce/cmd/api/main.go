package main

import (
	"net"
	"strings"

	"github.com/labstack/echo/v4"
	sensor "github.com/snipep/Assessment/Consumer-Service/Proto"
	"github.com/snipep/Assessment/Consumer-Service/internal/db"
	"github.com/snipep/Assessment/Consumer-Service/internal/env"
	"github.com/snipep/Assessment/Consumer-Service/internal/store"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	migratepkg "github.com/snipep/Assessment/Consumer-Service/cmd/migrate"
)

// @title Sensor Data Consumer Service API
// @version 1.0
// @description This is the API for the Consumer Service, which handles sensor data.
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and a JWT token.
func main() {
	// ---- Logger Setup ----
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	cfg := Config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "user:password@tcp(127.0.0.1:3306)/sensordata?charset=utf8mb4&parseTime=true&loc=UTC"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 25),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 25),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	// Sanitize env-provided values (Makefile include .env can leave literal quotes)
	cfg.addr = strings.Trim(cfg.addr, " \t\n\r\f\v\"'")
	dsn := strings.Trim(cfg.db.addr, " \t\n\r\f\v\"'")

	db, err := db.New(
		dsn,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run embedded migrations
	if err := migratepkg.Run(db); err != nil {
		logger.Fatalf("Failed to run migrations: %v", err)
	}

	logger.Info("database connection established")
	store := store.NewStorage(db)

	app := Application{
		config: cfg,
		store:  store,
		logger: logger,
	}

	// start worker pool for incoming data
	workers := env.GetInt("INGEST_WORKERS", 4)
	buffer := env.GetInt("INGEST_BUFFER", 1024)
	app.InitIngest(workers, buffer)

	// ---- gRPC Server Setup ----
	go startGRPCServer(&app)

	// ---- Start the HTTP server ----
	router := echo.New()
	RegisterMiddlewares(router)

	if err := app.run(router); err != nil {
		logger.Fatalf("Failed to run the application: %v", err)
	}

}

// startGRPCServer initializes and runs the gRPC server.
func startGRPCServer(app *Application) {
	// Support both GRPC_PORT and legacy GPRC_PORT
	port := env.GetString("GRPC_PORT", "")
	if port == "" {
		port = env.GetString("GPRC_PORT", ":50051")
	}
	p := strings.Trim(port, " \t\n\r\f\v\"'")
	if p != "" && p[0] != ':' {
		p = ":" + p
	}

	lis, err := net.Listen("tcp", p)
	if err != nil {
		app.logger.Fatalf("Failed to listen for gRPC: %v", err)
	}

	s := grpc.NewServer()
	grpcServer := NewGRPCServer(app)
	sensor.RegisterSensorServiceServer(s, grpcServer)

	app.logger.Infof("Starting gRPC server on %s", p)
	if err := s.Serve(lis); err != nil {
		app.logger.Fatalf("Failed to serve gRPC: %v", err)
	}
}
