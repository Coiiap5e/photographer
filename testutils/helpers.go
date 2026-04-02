package testutils

import (
	"context"
	"log/slog"
)

var testDB *TestDB

func CreateTestDB(ctx context.Context, logger *slog.Logger) (*TestDB, error) {
	var err error

	testDB, err = SetupTestDB(ctx)
	if err != nil {
		logger.Error("failed to setup test database", "error", err)
		return nil, err
	}

	return testDB, nil
}
