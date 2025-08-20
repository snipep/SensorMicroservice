package store

import (
	"context"
	"database/sql"
)

type SensorStore struct {
	db *sql.DB
}

type SensorData struct {
	SensorType  string  `json:"sensor_type"`
	SensorValue float32 `json:"sensor_value"`
	ID1         string  `json:"id1"`
	ID2         int32   `json:"id2"`
	Timestamp   string  `json:"timestamp"`
}

func (s SensorStore) InsertSensorData(ctx context.Context, data SensorData) error {
	query1 := `
		INSERT IGNORE INTO sensors (id1, sensor_type)
		VALUES (?, ?);
	`
	_, err := s.db.ExecContext(ctx, query1, data.ID1, data.SensorType)
	if err != nil {
		return err
	}

	query2 := `INSERT INTO sensor_readings (id1, id2, sensor_value, timestamp) VALUES (?, ?, ?, ?)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err = s.db.ExecContext(
		ctx,
		query2,
		data.ID1,
		data.ID2,
		data.SensorValue,
		data.Timestamp,
	)
	if err != nil {
		return err
	}

	return nil
}
