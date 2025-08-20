package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	sensor "github.com/snipep/Assessment/Consumer-Service/Proto"
	"github.com/snipep/Assessment/Consumer-Service/internal/db"
	"github.com/snipep/Assessment/Consumer-Service/internal/env"
	"github.com/snipep/Assessment/Consumer-Service/internal/store"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	// ---- Load environment variables from .env file ----
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	} 
	gRPCPort := os.Getenv("GPRC_PORT")
	if gRPCPort == "" {
		gRPCPort = ":50051" 	// Default port if not set in .env
	}

	cfg := Config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:	env.GetString("DB_ADDR", "user:password@tcp(127.0.0.1:3306)/sensordata?charset=utf8mb4&parseTime=True&loc=Local"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 25),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 25),
			maxIdleTime: env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}
	// ---- Logger Setup ----
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	db, err := db.New(
		cfg.db.addr, 
		cfg.db.maxOpenConns, 
		cfg.db.maxIdleConns, 
		cfg.db.maxIdleTime)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	logger.Info("database connection established")
	store := store.NewStorage(db)
	
	app := Application{
		config: cfg,
		store: store,
		logger: logger,
	}

	// ---- gRPC Server Setup ----
	lis, err := net.Listen("tcp", gRPCPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}


	// Create a new gRPC server instance.
	s := grpc.NewServer()
	// Register our server implementation with the gRPC server.
	sensor.RegisterSensorServiceServer(s, NewServer(store))
	log.Println("Server listening at", lis.Addr())
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
	// ---- Start the HTTP server ----
	router := echo.New()
	if err := app.run(router); err != nil {
		logger.Fatalf("Failed to run the application: %v", err)
	}
}
