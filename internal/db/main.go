package db

import (
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

func (s *Storage) Ping() error {
	return s.pg.Ping()
}
