package database

import (
	"context"
	"fmt"
	"time"

	"github.com/Coiiap5e/photographer/internal/config"
	"github.com/Coiiap5e/photographer/internal/errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func NewClient(ctx context.Context, db config.DbConfig, poolBuilder config.PoolConfig) (*DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		db.Username, db.Password, db.Host, db.Port, db.Database)

	return newFromConnStr(ctx, connStr, poolBuilder)
}

func newFromConnStr(ctx context.Context, connStr string, poolBuilder config.PoolConfig) (*DB, error) {
	configPool, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrCodeDBConnection, "failed to parse connection")
	}

	configPool.MaxConns = int32(poolBuilder.MaxOpenConns)
	configPool.MinConns = int32(poolBuilder.MaxIdleConns)
	configPool.MaxConnLifetime = poolBuilder.MaxConnLifetime
	configPool.MaxConnIdleTime = poolBuilder.MaxConnIdleTime

	pool, err := pgxpool.NewWithConfig(ctx, configPool)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrCodeDBConnection,
			"failed to connect to database")
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		return nil, errors.Wrap(err, errors.ErrCodeDBConnection,
			"error pinging database")
	}

	return &DB{Pool: pool}, nil
}

func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}
