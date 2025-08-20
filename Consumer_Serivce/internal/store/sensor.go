package store

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type SensorStore struct {
	db *sql.DB
}

type SensorData struct {
	SensorType  string  `json:"sensor_type"`
	SensorValue float32 `json:"sensor_value"`
	ID1         string  `json:"id1"`
	ID2         int32   `json:"id2"`
	Timestamp   time.Time  `json:"timestamp"`
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


func (s *SensorStore) GetSensorByIDs(id1 string, id2 int32) (*SensorData, error) {
	// --- FIX: Use a JOIN to get sensor_type from the 'sensors' table ---
	query := `
        SELECT r.id1, r.id2, r.sensor_value, r.timestamp, s.sensor_type
        FROM sensor_readings r
        INNER JOIN sensors s ON r.id1 = s.id1
        WHERE r.id1 = ? AND r.id2 = ?`

	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeoutDuration)
	defer cancel()
	row := s.db.QueryRowContext(ctx, query, id1, id2)

	sensor := &SensorData{}
	// --- FIX: Add sensor.SensorType to the Scan arguments ---
	err := row.Scan(
		&sensor.ID1,
		&sensor.ID2,
		&sensor.SensorValue,
		&sensor.Timestamp,
		&sensor.SensorType,
	)
	if err != nil {
		return nil, err
	}
	return sensor, nil
}


// --- NEW: Function to get all sensor readings within a time range ---
func (s *SensorStore) GetSensorHistory(startTime, endTime time.Time) ([]SensorData, error) {
	query := `
        SELECT r.id1, r.id2, r.sensor_value, r.timestamp, s.sensor_type
        FROM sensor_readings r
        INNER JOIN sensors s ON r.id1 = s.id1
        WHERE r.timestamp BETWEEN ? AND ?
        ORDER BY r.timestamp ASC`

	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sensors []SensorData
	for rows.Next() {
		var sensor SensorData
		// The Scan function will correctly parse the database's DATETIME format into the time.Time field.
		err := rows.Scan(
			&sensor.ID1,
			&sensor.ID2,
			&sensor.SensorValue,
			&sensor.Timestamp,
			&sensor.SensorType,
		)
		if err != nil {
			return nil, err
		}
		sensors = append(sensors, sensor)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	log.Print("sensors:", sensors)
	return sensors, nil
}

// --- NEW: Function to get sensor readings for specific IDs within a time range ---
func (s *SensorStore) GetSensorHistoryByIDs(id1 string, id2 int32, startTime, endTime time.Time) ([]SensorData, error) {
	query := `
        SELECT r.id1, r.id2, r.sensor_value, r.timestamp, s.sensor_type
        FROM sensor_readings r
        INNER JOIN sensors s ON r.id1 = s.id1
        WHERE r.id1 = ? 
          AND r.id2 = ? 
          AND r.timestamp BETWEEN ? AND ?
        ORDER BY r.timestamp ASC`

	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, id1, id2, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sensors []SensorData
	for rows.Next() {
		var sensor SensorData
		err := rows.Scan(
			&sensor.ID1,
			&sensor.ID2,
			&sensor.SensorValue,
			&sensor.Timestamp,
			&sensor.SensorType,
		)
		if err != nil {
			return nil, err
		}
		sensors = append(sensors, sensor)
	}
	log.Print("sensors:", sensors)

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sensors, nil
}

func (s *SensorStore) DeleteSensorDataByIDs(id1 string, id2 int32) (int64, error) {
	// The query is a simple DELETE with a WHERE clause for the specific ID pair.
	query := `DELETE FROM sensor_readings WHERE id1 = ? AND id2 = ?`

	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeoutDuration)
	defer cancel()

	// Use ExecContext to execute the delete statement with the provided IDs.
	result, err := s.db.ExecContext(ctx, query, id1, id2)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}


func (s *SensorStore) DeleteSensorHistory(startTime, endTime time.Time) (int64, error) {
	query := `DELETE FROM sensor_readings WHERE timestamp BETWEEN ? AND ?`

	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeoutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, startTime, endTime)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (s *SensorStore) DeleteSensorHistoryByIDs(id1 string, id2 int32, startTime, endTime time.Time) (int64, error) {
	query := `DELETE FROM sensor_readings WHERE id1 = ? AND id2 = ? AND timestamp BETWEEN ? AND ?`

	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeoutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, id1, id2, startTime, endTime)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
