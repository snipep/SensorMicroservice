package main

import (
	"math/rand/v2"
	"time"

	sensor "github.com/snipep/Assessment/Streamer-Service/Proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// CreateSensorPayload generates a sample SensorData message.
func CreateSensorPayload(i int) *sensor.SensorDataPayload {
	sensorMap := map[string]string{
		"A": "Temperature",
		"B": "Humidity",
		"C": "Pressure",
		"D": "Light",
	}

	// Extract keys into a slice so we can pick randomly
	keys := make([]string, 0, len(sensorMap))
	for k := range sensorMap {
		keys = append(keys, k)
	}

	// Pick a random key
	randKey := keys[rand.IntN(len(keys))]

	// Create a new SensorData message
	data := &sensor.SensorDataPayload{
		SensorValue: rand.Float32() * 100.0,
		SensorType:  sensorMap[randKey], // value
		Id1:         randKey,            // key
		Id2:         int32(1000 + i),
		Timestamp:   timestamppb.New(time.Now()),
	}
	return data
}
