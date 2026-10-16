package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func connectDB(ctx context.Context) (*pgxpool.Pool, error) {
	if err := godotenv.Load("../.env"); err != nil {

		return nil, err
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL пуст — .env не прочитан")
	}
	pool, err := pgxpool.New(ctx, databaseURL)
	if err == nil {
		pingErr := pool.Ping(ctx)
		if pingErr == nil {
			return pool, nil
		}
		return nil, fmt.Errorf("не удалось достучаться, %w", pingErr)

	}
	return nil, fmt.Errorf("не нащупал строку для подключения из .env, %w", err)
}
