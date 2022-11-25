package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

func NewPgConfig(username, password, host, port, database string) *PgConfig {
	return &PgConfig{Username: username, Password: password, Host: host, Port: port, Database: database}
}

func NewClient(ctx context.Context, cfg *PgConfig) (*pgxpool.Pool, error) {
	
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}
	
	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}
	return pool, nil
}
