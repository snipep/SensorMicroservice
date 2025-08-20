package main

import (
	"context"
	"fmt"
	"log"
	"time"

	sensor "github.com/snipep/Assessment/Streamer-Service/Proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewClient creates and returns a new gRPC client and its connection.
func NewClient(addr string) (*grpc.ClientConn, sensor.SensorServiceClient, error) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("did not connect: %w", err)
	}

	// Create a new client stub.
	client := sensor.NewSensorServiceClient(conn)
	return conn, client, nil
}

// SendDataStream handles the client-side streaming logic.
func SendDataStream(client sensor.SensorServiceClient) error {
	log.Println("Starting to send sensor data stream...")

	stream, err := client.SendSensorData(context.Background())
	if err != nil {
		return fmt.Errorf("could not open stream: %w", err)
	}

	for i := 0; i < 10; i++ {
		// Create a new sensor data payload 
		payload := CreateSensorPayload(i)
		req := &sensor.SensorData{Data: payload}

		if err := stream.Send(req); err != nil {
			return fmt.Errorf("failed to send data point %d: %w", i, err)
		}
		log.Printf("Sent data point #%d with value: %.2f", i, payload.SensorValue)

		time.Sleep(500 * time.Millisecond)
	}

	// After sending all the data, close the stream and receive the server's response.
	res, err := stream.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("failed to receive response: %w", err)
	}

	log.Printf("Server Response: Status=%s, Message='%s'", res.GetStatus(), res.GetMessage())
	return nil
}