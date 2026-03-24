package connectdb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type PostgresDatabase struct {
	db      *sql.DB
	connStr string
}

func NewPostgresDatabase(connStr string) (*PostgresDatabase, error) {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return &PostgresDatabase{
		db:      db,
		connStr: connStr,
	}, nil
}

func (p *PostgresDatabase) GetDB() *sql.DB {
	return p.db
}

func (p *PostgresDatabase) Close() error {
	return p.db.Close()
}

func CheckDBConnection(db *sql.DB) error {
	return db.Ping()
}

func (p *PostgresDatabase) Ping() error {
	return p.db.Ping()
}

func (p *PostgresDatabase) Reconnect() error {

	newDB, err := sql.Open("postgres", p.connStr)
	if err != nil {
		return err
	}

	if err = newDB.Ping(); err != nil {
		return err
	}

	p.db = newDB
	return nil
}
