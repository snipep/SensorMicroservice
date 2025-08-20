package main

import (
	"fmt"
	"io"
	"log"
	"time"

	sensor "github.com/snipep/Assessment/Consumer-Service/Proto"
	"github.com/snipep/Assessment/Consumer-Service/internal/store"
)

type Server struct {
	sensor.UnimplementedSensorServiceServer
	Storage *store.Storage
}

func NewServer(storage *store.Storage) *Server {
	return &Server{
		Storage: storage,
	}
}

func (s *Server) SendSensorData(stream sensor.SensorService_SendSensorDataServer) error {
	log.Println("Client connected. Starting to receive stream.")
	var dataCount int32

	for {
		// Recv() blocks until a message is received or the stream is closed.
		// The received message type is now *sensor.SensorData
		req, err := stream.Recv()

		// If the stream is closed by the client, io.EOF is returned.
		if err == io.EOF {
			log.Printf("Finished receiving data. Total messages: %d", dataCount)
			// Send a final response back to the client and close the connection.
			return stream.SendAndClose(&sensor.SensorDataResponse{
				Status:  "Success",
				Message: fmt.Sprintf("Successfully processed %d data points.", dataCount),
			})
		}
		// Handle any other errors during reception.
		if err != nil {
			log.Printf("Error while receiving stream: %v", err)
			return err
		}

		payload := req.GetData()
		var timestamp time.Time
		if payload.GetTimestamp() == nil {
			timestamp = time.Now()
		} else {
			timestamp = payload.GetTimestamp().AsTime()
		}
		
		dbData := store.SensorData{
			SensorType:  payload.GetSensorType(),
			SensorValue: payload.GetSensorValue(),
			ID1:         payload.GetId1(),
			ID2:         payload.GetId2(),
			Timestamp:   timestamp.Format("2006-01-02 15:04:05"),
		}

		err = s.Storage.Sensor.InsertSensorData(stream.Context(), dbData)
		if err != nil {
			log.Printf("ERROR: Failed to insert sensor data: %v", err)
		} else {
			log.Printf("Successfully inserted: Type=%s, Value=%.2f, ID1=%s",
				dbData.SensorType, dbData.SensorValue, dbData.ID1)
		}
		dataCount++
	}
}
