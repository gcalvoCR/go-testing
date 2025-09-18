package main

import (
	"fmt"
	"testing"

	"github.com/gcalvocr/go-testing/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceGetUserWithExpect(t *testing.T) {
	mockRepo := new(mocks.MockUserRepositoryWithExpecter)

	mockRepo.EXPECT().GetUser(mock.Anything).Return("Alice", nil).Run(func(id int) {
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
