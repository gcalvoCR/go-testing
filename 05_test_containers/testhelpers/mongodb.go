package testhelpers

import (
	"context"
	"testcontainers-mongo-demo/config"
	"testing"

	"github.com/testcontainers/testcontainers-go/modules/mongodb"
)

// MongoDBTestHelper provides utilities for testing with MongoDB containers
type MongoDBTestHelper struct {
	Container *mongodb.MongoDBContainer
	Config    *config.DatabaseConfig
}

// SetupMongoDBContainer creates and configures a MongoDB test container
func SetupMongoDBContainer(ctx context.Context, t *testing.T) *MongoDBTestHelper {
	t.Helper()

	// Start MongoDB test container
	mongoContainer, err := mongodb.Run(ctx, "mongo:7")
	if err != nil {
		t.Fatal(err)
	}

	// Get connection string
	uri, err := mongoContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Create config
	cfg := &config.DatabaseConfig{
		URI:        uri,
		Database:   "testdb",
		Collection: "users",
	}

	return &MongoDBTestHelper{
		Container: mongoContainer,
		Config:    cfg,
		// Repo:      repo,
	}
}

// Cleanup terminates the container and closes the repository
func (h *MongoDBTestHelper) Cleanup(ctx context.Context, t *testing.T) {
	t.Helper()

	// Terminate container
	if err := h.Container.Terminate(ctx); err != nil {
		t.Logf("Failed to terminate container: %v", err)
	}
}
