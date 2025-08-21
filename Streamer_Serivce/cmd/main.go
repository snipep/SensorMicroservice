package main

import (

	"os"
	"time"
	"go.uber.org/zap"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	err := godotenv.Load(".env")
	if err != nil {
		logger.Errorf("Error loading .env file: %v", err)
	}

	serverAddr := os.Getenv("MICROSERVICE_B_ADDR")
	if serverAddr == "" {
		serverAddr = "localhost:50051" // Default address if not set in .env
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = ":8080"
	}

	// Establish a connection and get a new client.
	conn, gclient, err := NewClient(serverAddr)
	if err != nil {
		logger.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()
	logger.Info("Successfully connected to gRPC server.")
	
	
	
	// ---- Application State ----
	app := NewApplication(gclient, time.NewTicker(5*time.Second), logger)
	router := echo.New()
	// Register routes for the Application.
	app.RegiterRoutes(router)
	// Serve the HTTP server.
	
	// Start the data streaming in a separate goroutine.
	go app.SendDataStream()
	
	if err := router.Start(httpPort); err != nil {
		logger.Fatalf("Failed to start HTTP server: %v", err)
	}
}
