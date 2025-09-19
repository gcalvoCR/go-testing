package main

import (
	"context"
	"testcontainers-mongo-demo/testhelpers"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_InitializeDatabase(t *testing.T) {
	ctx := context.Background()

	helper := testhelpers.SetupMongoDBContainer(ctx, t)
	defer helper.Cleanup(ctx, t)

	// Set environment variables for test
	t.Setenv("MONGO_URI", helper.Config.URI)
	t.Setenv("MONGO_DATABASE", helper.Config.Database)
	t.Setenv("MONGO_COLLECTION", helper.Config.Collection)

	// Test InitializeDatabase
	repo, err := InitializeDataBase()
	require.NoError(t, err)

	defer func() {
		if err := repo.Close(ctx); err != nil {
			t.Logf("Failed to close repository: %v", err)
		}
	}()

	// Verify repository works
	err = repo.Ping(ctx)
	require.NoError(t, err)
}
