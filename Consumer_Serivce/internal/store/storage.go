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
	Sensor interface {
		InsertSensorData(ctx context.Context, data SensorData) error
		GetSensorByIDs(id1 string, id2 int32, limit, offset int) (*SensorData, error)
		GetSensorHistory(startTime, endTime time.Time, limit, offset int) ([]SensorData, error)
		GetSensorHistoryByIDs(id1 string, id2 int32, startTime, endTime time.Time, limit, offset int) ([]SensorData, error)
		DeleteSensorHistoryByIDs(id1 string, id2 int32, startTime, endTime time.Time) (int64, error)
		DeleteSensorHistory(startTime, endTime time.Time) (int64, error)
		DeleteSensorDataByIDs(id1 string, id2 int32) (int64, error)
		EditSensorDataByID(id1 string, id2 int32, newValue float32) (int64, error)
		EditSensorHistory(startTime, endTime time.Time, newValue float32) (int64, error)
		EditSensorHistoryByIDs(id1 string, id2 int32, startTime, endTime time.Time, newValue float32) (int64, error)
	}
	User interface {
		CreateUser(ctx context.Context, name, email, passwordHash string) (int64, error) 
		GetUserByEmail(ctx context.Context, email string) (*User, error)
	}
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Sensor: &SensorStore{db: db},
		User:   &UserStore{db: db},
	}
}
