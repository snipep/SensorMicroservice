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

func NewClient(addr string) (*grpc.ClientConn, sensor.SensorServiceClient, error) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("did not connect: %w", err)
	}

	client := sensor.NewSensorServiceClient(conn)
	return conn, client, nil
}

func (app *Application) SendDataStream() {
	log.Println("Starting background data streaming goroutine...")
	var stream sensor.SensorService_SendSensorDataClient
	var err error
	i := 0 

	for {
		// If the stream is nil,Create a new one.
		if stream == nil {
			stream, err = app.SensorClient.SendSensorData(context.Background())
			if err != nil {
				app.logger.Errorf("Could not open stream, will retry in 5 seconds: %v", err)
				time.Sleep(5 * time.Second)
				continue // Reconnect
			}
			log.Println("Successfully established new gRPC stream.")
		}

		<-app.ticker.C

		payload := app.CreateSensorPayload(i)
		req := &sensor.SensorData{Data: payload}

		// Send the data. If it fails, we set the stream to nil to force a reconnect on the next loop iteration.
		if err := stream.Send(req); err != nil {
			app.logger.Errorf("Failed to send data point, will attempt to reconnect: %v", err)
			stream = nil // Mark stream as broken
			continue
		}

		log.Printf("Sent data point #%d with value: %.2f", i, payload.SensorValue)
		i++
	}
}
