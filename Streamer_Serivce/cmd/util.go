package main

import (
	"math/rand/v2"
	"time"

	sensor "github.com/snipep/Assessment/Streamer-Service/Proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// CreateSensorPayload generates a sample SensorData message for the application's fixed sensor.
func (app *Application) CreateSensorPayload(i int) *sensor.SensorDataPayload {
	// generate a pseudo-random value for the fixed sensor
	val := rand.Float32() * 100.0
	data := &sensor.SensorDataPayload{
		SensorValue: val,
		SensorType:  app.FixedSensorType,
		Id1:         app.FixedSensorID1,
		Id2:         int32(1000 + i),
		Timestamp:   timestamppb.New(time.Now()),
	}
	return data
}
