package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"main.go/config"
)

func NewPostgresPool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", cfg.Database.User, cfg.Database.Pasword, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name, cfg.Database.SSLMode)

	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
