package main

import (
	"math/rand/v2"
	"time"

	sensor "github.com/snipep/Assessment/Streamer-Service/Proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (app *Application) CreateSensorPayload(i int) *sensor.SensorDataPayload {
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
