package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type User struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) CreateUser(ctx context.Context, name, email, passwordHash string) (int64, error) {
	query := `INSERT INTO users (name, email, password_hash, created_at) VALUES (?, ?, ?, ?)`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	createdAt := time.Now().UTC()
	res, err := s.db.ExecContext(ctx, query, name, email, passwordHash, createdAt)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *UserStore) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, name, email, password_hash, created_at FROM users WHERE email = ? LIMIT 1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	row := s.db.QueryRowContext(ctx, query, email)
	var u User
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return &u, nil
}
