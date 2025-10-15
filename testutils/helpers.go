package testutils

import (
	"context"
	"log"
)

var testDB *TestDB

func CreateTestDB(ctx context.Context) (*TestDB, error) {
	var err error

	testDB, err = SetupTestDB(ctx)
	if err != nil {
		log.Printf("Failed to setup test database: %v", err)
		return nil, err
	}

	return testDB, nil
}
