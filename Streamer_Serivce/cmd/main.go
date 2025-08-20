package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	serverAddr := os.Getenv("MICROSERVICE_B_ADDR")
	if serverAddr == "" {
		serverAddr = "localhost:50051" // Default address if not set in .env
	}

	// Establish a connection and get a new client.
	conn, client, err := NewClient(serverAddr)
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	log.Println("Successfully connected to gRPC server.")

	// Call the function to send a stream of sensor data.
	if err := SendDataStream(client); err != nil {
		log.Fatalf("Failed to send data stream: %v", err)
	}

	log.Println("Client execution finished.")
}
