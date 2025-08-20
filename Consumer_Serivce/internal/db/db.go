package db

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// New creates a new database connection pool for MySQL.
func New(dsn string, maxOpenConns int, maxIdleConns int, maxIdleTime string) (*sql.DB, error) {
	// Use "mysql" as the driver name.
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)

	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Ping the database to verify the connection is alive.
	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}