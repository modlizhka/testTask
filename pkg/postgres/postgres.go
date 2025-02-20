package postgres

import (
	"context"
	"fmt"
	"user-service/config"

	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type Pool struct {
	*pgxpool.Pool
}

func NewClient(ctx context.Context, cfg config.DataBase) (pool Pool, err error) {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	p, err := pgxpool.Connect(ctx, connStr)
	if err != nil {
		return Pool{nil}, err
	}
	return Pool{p}, nil
}
