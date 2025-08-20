package store

import "database/sql"

type SensorStore struct{
	db *sql.DB
}