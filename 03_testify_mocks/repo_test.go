package main

import (
	"fmt"
	"testing"

	"github.com/gcalvocr/go-testing/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceGetUser(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)

	// Tell the mock: "If GetUser(42) is called, return 'Alice', nil"
	// mockRepo.On("GetUser", 42).Return("Alice", nil)

	mockRepo.On("GetUser", mock.Anything).Return("Bob", nil).Run(func(args mock.Arguments) {
		id := args.Int(0)
		fmt.Println("GetUser called with ID:", id)
	})

	// Inject mock into your service
	service := NewUserService(mockRepo)

	user, err := service.GetUser(42)

	assert.NoError(t, err)
	assert.Equal(t, "Alice", user)

	// Verify that expectations were met
	mockRepo.AssertExpectations(t)
}
