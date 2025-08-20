package main

import (
	"math/rand/v2"
	"time"

	sensor "github.com/snipep/Assessment/Streamer-Service/Proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// CreateSensorPayload generates a sample SensorData message.
func CreateSensorPayload(i int) *sensor.SensorDataPayload {
	sensorType := []string{"Temperature", "Humidity", "Pressure", "Light"}

	id1Runes := make([]rune, 4)
	for i := range id1Runes {
		id1Runes[i] = rune('A' + rand.IntN(26)) 
	}
	// Create a new SensorData message.
	data := &sensor.SensorDataPayload{
		SensorValue: rand.Float32() * 100.0,
		SensorType:  sensorType[rand.IntN(len(sensorType))], 
		Id1:         string(id1Runes),
		Id2:         int32(1000 + i),
		Timestamp: timestamppb.New(time.Now()),
	}
	return data
}