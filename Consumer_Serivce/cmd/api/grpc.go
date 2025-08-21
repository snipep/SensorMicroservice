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
	Application *Application
}

func NewGRPCServer(app *Application) *Server {
	return &Server{
		Application: app,
	}
}

func (s *Server) SendSensorData(stream sensor.SensorService_SendSensorDataServer) error {
	log.Println("Client connected. Starting to receive stream.")
	var dataCount int32

	for {
		// Recv() blocks until a message is received or the stream is closed.
		req, err := stream.Recv()

		if err == io.EOF {
			log.Printf("Finished receiving data. Total messages: %d", dataCount)
			return stream.SendAndClose(&sensor.SensorDataResponse{
				Status:  "Success",
				Message: fmt.Sprintf("Successfully processed %d data points.", dataCount),
			})
		}
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
			Timestamp:   timestamp,
		}
		
		// Send to ingest channel (unbuffered). This will block until a worker picks it up.
		s.Application.ingestCh <- dbData
		log.Printf("Queued for insert: Type=%s, Value=%.2f, ID1=%s", dbData.SensorType, dbData.SensorValue, dbData.ID1)
		dataCount++
	}
}
