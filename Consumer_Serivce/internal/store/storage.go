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
		GetSensorByIDs(id1 string, id2 int32) (*SensorData, error)
		GetSensorHistory(startTime, endTime time.Time) ([]SensorData, error)
		GetSensorHistoryByIDs(id1 string, id2 int32, startTime, endTime time.Time) ([]SensorData, error)
	}
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Sensor: &SensorStore{db: db},
	}
}