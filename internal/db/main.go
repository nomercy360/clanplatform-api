package db

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Init initializes a new database connection.
func Init(connStr string) (Storage, error) {
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return Storage{}, err
	}

	db.SetConnMaxLifetime(60)
	db.SetConnMaxIdleTime(30)
	db.SetMaxIdleConns(3)
	db.SetMaxOpenConns(15)

	log.Println("Database connection established")

	return Storage{pg: db}, nil
}

// Close closes the database connection.
func (s *Storage) Close() {
	if s.pg != nil {
		err := s.pg.Close()
		if err != nil {
			log.Println("Failed to close database connection:", err)
		} else {
			log.Println("Database connection closed")
		}
	}
}

func IsDuplicationError(err error) bool {
	if err == nil {
		return false
	}

	var pqErr *pq.Error
	ok := errors.As(err, &pqErr)
	return ok && pqErr.Code == "23505"
}

func IsForeignKeyViolationError(err error) bool {
	if err == nil {
		return false
	}

	var pqErr *pq.Error
	ok := errors.As(err, &pqErr)
	return ok && pqErr.Code == "23503"
}

// IsNoRowsError check for no rows found error
func IsNoRowsError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func (s *Storage) Ping() error {
	err := s.pg.Ping()

	return err
}
