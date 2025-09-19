package demo

import (
	"context"
	"fmt"
	"testcontainers-mongo-demo/models"
	"testcontainers-mongo-demo/repository"
)

// RunUserDemo performs a demonstration of user repository operations
func RunUserDemo(ctx context.Context, repo repository.UserRepository) error {
	fmt.Println("Running demo operations...")

	// Create a user
	user := &models.User{
		Name: "John",
		Age:  30,
	}
	err := repo.Create(ctx, user)
	if err != nil {
		return err
	}
	fmt.Println("User created with ID:", user.ID)

	// Get all users
	users, err := repo.GetAll(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("Total users: %d\n", len(users))

	return nil
}
