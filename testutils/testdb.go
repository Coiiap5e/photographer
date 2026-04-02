package testutils

import (
	"context"
	"fmt"
	"time"

	"github.com/Coiiap5e/photographer/internal/errors"
	"github.com/Coiiap5e/photographer/internal/infrastructure/database"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDB struct {
	Container testcontainers.Container
	Pool      *pgxpool.Pool
}

func SetupTestDB(ctx context.Context) (*TestDB, error) {

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "photographer_test",
			"POSTGRES_USER":     "test_user",
			"POSTGRES_PASSWORD": "test_password",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get host: %w", err)
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, fmt.Errorf("failed to get port: %w", err)
	}

	time.Sleep(3 * time.Second)

	dsn := fmt.Sprintf("postgresql://test_user:test_password@%s:%s/photographer_test?sslmode=disable",
		host, port.Port())

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	testDB := &TestDB{
		Container: container,
		Pool:      pool,
	}

	if err := testDB.InitSchema(ctx); err != nil {
		return nil, fmt.Errorf("failed to init schema: %w", err)
	}

	return testDB, nil
}

func (tdb *TestDB) InitSchema(ctx context.Context) error {

	_, err := tdb.Pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS clients (
			id SERIAL PRIMARY KEY,
			first_name VARCHAR(100) NOT NULL,
			last_name VARCHAR(100) NOT NULL,
			phone VARCHAR(100),
			social_network_url VARCHAR(500),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return errors.Wrap(err, errors.ErrCodeClientCreate, "failed to create clients table")
	}

	_, err = tdb.Pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS shoots (
			id SERIAL PRIMARY KEY,
			date DATE NOT NULL,
			start_time TIME NOT NULL,
			end_time TIME NOT NULL,
			shoot_price DECIMAL(10,0),
			location VARCHAR(255),
			shoot_type VARCHAR(100),
			notes TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return errors.Wrap(err, errors.ErrCodeShootCreate, "failed to create shoots table")
	}

	_, err = tdb.Pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS shoot_clients (
			shoot_id INTEGER REFERENCES shoots(id) ON DELETE CASCADE,
			client_id INTEGER REFERENCES clients(id) ON DELETE CASCADE,
			is_main_client BOOLEAN DEFAULT false,
			relationship_type VARCHAR(50) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (shoot_id, client_id)
		)
	`)

	if err != nil {
		return errors.Wrap(err, errors.ErrCodeShootCreate, "failed to create shoot client table")
	}

	return nil
}

func (tdb *TestDB) Cleanup(ctx context.Context) error {

	if tdb.Pool != nil {
		tdb.Pool.Close()
	}

	if tdb.Container != nil {
		if err := tdb.Container.Terminate(ctx); err != nil {
			return fmt.Errorf("failed to terminate container: %w", err)
		}
	}

	return nil
}

func (tdb *TestDB) CleanTables(ctx context.Context) error {
	_, err := tdb.Pool.Exec(ctx, "TRUNCATE TABLE shoots, clients RESTART IDENTITY CASCADE")
	if err != nil {
		return errors.Wrap(err, errors.ErrCodeDBDelete, "failed to clean tables")
	}
	return nil
}

func (tdb *TestDB) GetDB() *database.DB {
	return &database.DB{
		Pool: tdb.Pool,
	}
}
