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
		DeleteSensorHistoryByIDs(id1 string, id2 int32, startTime, endTime time.Time) (int64, error)
		DeleteSensorHistory(startTime, endTime time.Time) (int64, error)
		DeleteSensorDataByIDs(id1 string, id2 int32) (int64, error)
		EditSensorDataByID(id1 string, id2 int32, newValue float32) (int64, error)
		EditSensorHistory(startTime, endTime time.Time, newValue float32) (int64, error)
		EditSensorHistoryByIDs(id1 string, id2 int32, startTime, endTime time.Time, newValue float32) (int64, error)
	}
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Sensor: &SensorStore{db: db},
}
}