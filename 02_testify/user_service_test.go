package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockUserDatabase is a mock implementation of UserDatabase
type MockUserDatabase struct {
	mock.Mock
}

func (m *MockUserDatabase) GetUserByID(id int) (*User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserDatabase) CreateUser(user *User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserDatabase) UpdateUser(user *User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserDatabase) DeleteUser(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestUserService_GetUser(t *testing.T) {
	mockDB := new(MockUserDatabase)
	service := NewUserService(mockDB)

	t.Run("successful user retrieval", func(t *testing.T) {
		expectedUser := &User{ID: 1, Name: "John Doe", Email: "john@example.com"}

		// Setup mock expectations
		mockDB.On("GetUserByID", 1).Return(expectedUser, nil).Once()

		// Execute
		user, err := service.GetUser(1)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, expectedUser.ID, user.ID)
		assert.Equal(t, expectedUser.Name, user.Name)
		assert.Equal(t, expectedUser.Email, user.Email)

		// Verify mock expectations
		mockDB.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		// Setup mock expectations
		mockDB.On("GetUserByID", 999).Return(nil, errors.New("user not found")).Once()

		// Execute
		user, err := service.GetUser(999)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "failed to get user")

		// Verify mock expectations
		mockDB.AssertExpectations(t)
	})

	t.Run("invalid user ID", func(t *testing.T) {
		// Execute - should not call database
		user, err := service.GetUser(0)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "invalid user ID", err.Error())

		// Verify database was not called
		mockDB.AssertNotCalled(t, "GetUserByID", mock.Anything)
	})
}

func TestUserService_CreateUser(t *testing.T) {
	mockDB := new(MockUserDatabase)
	service := NewUserService(mockDB)

	t.Run("successful user creation", func(t *testing.T) {
		// Setup mock expectations
		mockDB.On("CreateUser", mock.MatchedBy(func(user *User) bool {
			return user.Name == "Jane Doe" && user.Email == "jane@example.com"
		})).Return(nil).Once()

		// Execute
		user, err := service.CreateUser("Jane Doe", "jane@example.com")

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "Jane Doe", user.Name)
		assert.Equal(t, "jane@example.com", user.Email)

		// Verify mock expectations
		mockDB.AssertExpectations(t)
	})

	t.Run("validation error - empty name", func(t *testing.T) {
		// Execute
		user, err := service.CreateUser("", "jane@example.com")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "name is required", err.Error())

		// Verify database was not called
		mockDB.AssertNotCalled(t, "CreateUser", mock.Anything)
	})

	t.Run("validation error - empty email", func(t *testing.T) {
		// Execute
		user, err := service.CreateUser("Jane Doe", "")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "email is required", err.Error())

		// Verify database was not called
		mockDB.AssertNotCalled(t, "CreateUser", mock.Anything)
	})

	t.Run("database error", func(t *testing.T) {
		// Setup mock expectations
		mockDB.On("CreateUser", mock.Anything).Return(errors.New("database connection failed")).Once()

		// Execute
		user, err := service.CreateUser("Jane Doe", "jane@example.com")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "failed to create user")

		// Verify mock expectations
		mockDB.AssertExpectations(t)
	})
}

func TestUserService_UpdateUserEmail(t *testing.T) {
	mockDB := new(MockUserDatabase)
	service := NewUserService(mockDB)

	t.Run("successful email update", func(t *testing.T) {
		existingUser := &User{ID: 1, Name: "John Doe", Email: "john@example.com"}

		// Setup mock expectations
		mockDB.On("GetUserByID", 1).Return(existingUser, nil).Once()
		mockDB.On("UpdateUser", mock.MatchedBy(func(user *User) bool {
			return user.ID == 1 && user.Email == "john.doe@example.com"
		})).Return(nil).Once()

		// Execute
		err := service.UpdateUserEmail(1, "john.doe@example.com")

		// Assert
		assert.NoError(t, err)

		// Verify mock expectations
		mockDB.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		// Setup mock expectations
		mockDB.On("GetUserByID", 999).Return(nil, errors.New("user not found")).Once()

		// Execute
		err := service.UpdateUserEmail(999, "new@example.com")

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get user for update")

		// Verify mock expectations
		mockDB.AssertExpectations(t)
	})
}
