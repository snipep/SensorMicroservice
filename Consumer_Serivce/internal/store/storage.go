package store

import (
	"context"
	"database/sql"
	"time"
)

var (
	QueryTimeoutDuration = time.Second * 3
)

type Storage struct {
	Sensor interface{
		InsertSensorData(ctx context.Context, data SensorData) error
	}
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Sensor: &SensorStore{db: db},
	}
}