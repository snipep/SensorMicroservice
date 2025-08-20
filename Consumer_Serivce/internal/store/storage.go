package store

import "database/sql"

type Storage struct {
	Sensor interface{

	}
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Sensor: SensorStore{db: db},
	}
}