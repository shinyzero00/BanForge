package storage

import (
	"database/sql"

	"github.com/d3m0k1d/BanForge/internal/logger"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	logger *logger.Logger
	db     *sql.DB
}

func NewDB(path string) (*DB, error) {
	db, err := sql.Open("sqlite3", "/var/lib/banforge/storage.db")
	if err != nil {
		return nil, err
	}
	return &DB{
		logger: logger.New(false),
		db:     db,
	}, nil
}

func (d *DB) Close() error {
	d.logger.Info("Closing database connection")
	err := d.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) CreateTable() error {
	_, err := d.db.Exec(CreateTables)
	if err != nil {
		return err
	}
	d.logger.Info("Created tables")
	return nil
}
