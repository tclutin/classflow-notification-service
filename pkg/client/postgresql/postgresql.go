package postgresql

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"os"
	"time"
)

const (
	maxRetries = 5
)

type Client interface {
}

func NewPool(ctx context.Context, dsn string) *pgxpool.Pool {
	for i := 0; i < maxRetries; i++ {
		pool, err := pgxpool.New(ctx, dsn)
		if err != nil {
			slog.Error("failed to connect to the database", "retry_count", i+1, "error", err)
			time.Sleep(3 * time.Second)
			continue
		}

		if err = pool.Ping(ctx); err != nil {
			slog.Error("failed to ping database, retrying...", "retry_count", i+1, "error", err)
			time.Sleep(3 * time.Second)
			continue
		}

		return pool
	}

	slog.Error("max retries reached, unable to connect to the database")
	os.Exit(1)
	return nil
}
