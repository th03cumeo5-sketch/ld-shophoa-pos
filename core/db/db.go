package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/alexbrainman/odbc"
	_ "github.com/microsoft/go-mssqldb"
	"shophoa/core/config"
)

// Connect initializes and validates a connection pool to SQL Server
func Connect(cfg *config.Config) (*sql.DB, error) {
	connStr := cfg.ConnectionString()
	driverName := cfg.Driver
	if driverName == "" {
		driverName = "odbc"
	}

	db, err := sql.Open(driverName, connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open sql connection with driver %s: %w", driverName, err)
	}

	// Optimize connection pool for POS and back-office usage
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(10 * time.Minute)
	db.SetConnMaxIdleTime(3 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping failed to sql server %s (driver %s): %w", cfg.DBServer, driverName, err)
	}

	return db, nil
}
