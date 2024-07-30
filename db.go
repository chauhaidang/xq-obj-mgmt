package main

import (
	"database/sql"
	"fmt"
	"os"
)

type PGStorage struct {
	db *sql.DB
}

func NewPGStorage(cfg Config) (*PGStorage, error) {
	conn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.DBHost, 5432, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PGStorage{
		db: db,
	}, nil
}

func (s *PGStorage) Init() (*sql.DB, error) {
	dat, err := os.ReadFile("./x-obj-mgmt.sql")
	if err != nil {
		return nil, err
	}
	query := string(dat)
	_, err = s.db.Exec(query)
	if err != nil {
		return nil, err
	}
	return s.db, nil
}
