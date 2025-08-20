package main

import (
	"fmt"
	"io"
	"log"

	sensor "github.com/snipep/Assessment/Consumer-Service/Proto"
)

type Server struct {
	sensor.UnimplementedSensorServiceServer
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

		// Process the received data by getting the payload.
		data := req.GetData()
		log.Printf("Received: Type=%s, Value=%.2f, ID1=%s, ID2=%d",data.GetSensorType(), data.GetSensorValue(), data.GetId1(), data.GetId2())
		dataCount++
	}
}
