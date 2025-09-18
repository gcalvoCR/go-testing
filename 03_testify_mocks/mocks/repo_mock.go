package mocks

import "github.com/stretchr/testify/mock"

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUser(id int) (string, error) {
	args := m.Called(id) // track that it was called with `id`
	return args.String(0), args.Error(1)
}
