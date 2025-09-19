package repository

import (
	"context"
	"testcontainers-mongo-demo/models"
	"testcontainers-mongo-demo/testhelpers"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_MongoUserRepository(t *testing.T) {
	ctx := context.Background()

	helper := testhelpers.SetupMongoDBContainer(ctx, t)
	defer helper.Cleanup(ctx, t)

	repo, err := NewMongoUserRepository(helper.Config)
	require.NoError(t, err)

	// Test Create
	user := &models.User{
		Name: "Test",
		Age:  25,
	}
	err = repo.Create(ctx, user)
	require.NoError(t, err)
	require.NotEmpty(t, user.ID, "User ID should be set")

	// Test GetByID
	retrieved, err := repo.GetByID(ctx, user.ID)
	require.NoError(t, err)
	require.Equal(t, "Test", retrieved.Name)

	// Test Update
	user.Age = 26
	err = repo.Update(ctx, user)
	require.NoError(t, err)
	updated, err := repo.GetByID(ctx, user.ID)
	require.NoError(t, err)
	require.Equal(t, 26, updated.Age)

	// Test GetAll
	users, err := repo.GetAll(ctx)
	require.NoError(t, err)
	require.Len(t, users, 1)

	// Test Delete
	err = repo.Delete(ctx, user.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, user.ID)
	require.Error(t, err, "Expected error when getting deleted user")
}
